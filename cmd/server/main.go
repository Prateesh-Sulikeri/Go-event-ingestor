package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/api"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/auth"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/config"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/limiter"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	// Initialize Redis
	rdb := redis.NewClient(&redis.Options{Addr: cfg.REDIS_ADDR})
	lim := limiter.NewLimiter(rdb, cfg.RATE_LIMIT_BUCKET_SIZE, cfg.RATE_LIMIT_REFILL_RATE)

	// Initialize JWT
	jwtService := auth.NewJWTService(cfg)

	// Initialize DB store
	store, err := storage.NewEventStore(cfg)
	if err != nil {
		log.Fatal("db connect failed:", err)
	}
	eventHandler := api.NewEventHandler(store)

	// Gin setup
	r := gin.Default()

	// Public route: token generation
	r.POST("/auth/login", func(c *gin.Context) {
		var body struct{ ClientID string `json:"client_id"` }
		if err := c.BindJSON(&body); err != nil || body.ClientID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "client_id required"})
			return
		}
		token, err := jwtService.GenerateToken(body.ClientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	// Protected API group
	apiGroup := r.Group("/v1")
	apiGroup.Use(jwtService.Middleware())
	apiGroup.Use(lim.Middleware(func(c *gin.Context) string {
		val, _ := c.Get("client_id")
		return val.(string)
	}))

	// Rate-limited, authenticated routes
	apiGroup.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})

	apiGroup.POST("/events", eventHandler.Ingest)

	// Start service
	addr := ":8080"
	fmt.Println("Server running on", addr)
	log.Fatal(r.Run(addr))
}
