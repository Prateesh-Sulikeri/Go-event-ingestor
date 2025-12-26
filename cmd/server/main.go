package main

import (
	"fmt"
	"log"

	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/api"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/auth"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/config"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/limiter"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/metrics"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	// Core services
	rdb := redis.NewClient(&redis.Options{Addr: cfg.REDIS_ADDR})
	lim := limiter.NewLimiter(rdb, cfg.RATE_LIMIT_BUCKET_SIZE, cfg.RATE_LIMIT_REFILL_RATE)
	m := metrics.NewMetrics()
	jwt := auth.NewJWTService(cfg)

	store, err := storage.NewEventStore(cfg)
	if err != nil {
		log.Fatal("DB connect failed:", err)
	}
	eventHandler := api.NewEventHandler(store)

	// Gin router
	r := gin.Default()

	// Auth
	r.POST("/auth/login", func(c *gin.Context) {
		var body struct{ ClientID string `json:"client_id"` }
		if err := c.BindJSON(&body); err != nil || body.ClientID == "" {
			c.JSON(400, gin.H{"error": "client_id required"})
			return
		}
		token, _ := jwt.GenerateToken(body.ClientID)
		c.JSON(200, gin.H{"token": token})
	})

	// WebSocket for live metrics
	wsHandler := api.NewWebSocketHandler(m, lim, rdb)
	r.GET("/v1/ws", wsHandler.Handle)

	// Protected routes (ORDER MATTERS)
	v1 := r.Group("/v1")
	v1.Use(jwt.Middleware())

	// Metrics then limiter
	v1.Use(m.Middleware(func(c *gin.Context) string {
		id, _ := c.Get("client_id")
		return id.(string)
	}))
	v1.Use(lim.Middleware(func(c *gin.Context) string {
		id, _ := c.Get("client_id")
		return id.(string)
	}, m))

	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"msg": "pong"})
	})

	v1.POST("/events", eventHandler.Ingest)

	// Dynamic limiter config
	v1.POST("/config/limiter", func(c *gin.Context) {
		var b struct {
			Bucket int `json:"bucket"`
			Refill int `json:"refill"`
		}
		if err := c.BindJSON(&b); err != nil {
			c.JSON(400, gin.H{"error": "invalid payload"})
			return
		}
		lim.Update(b.Bucket, b.Refill)
		m.BucketSize = b.Bucket
		m.RefillRate = b.Refill
		c.JSON(200, gin.H{"status": "updated"})
	})

	// Dashboard UI
	r.Static("/static", "./dashboard/static")
	r.LoadHTMLFiles("dashboard/index.html")

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	fmt.Println("Server running on :8080")
	log.Fatal(r.Run(":8080"))
}
