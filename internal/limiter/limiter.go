package limiter

import (
	"context"
	_ "embed"
	"net/http"
	"time"

	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

//go:embed token_bucket.lua
var luaScript string

// Extract client ID from context

type Limiter struct {
	client *redis.Client
	max    int // bucket size
	refill int // tokens per second
}

func NewLimiter(client *redis.Client, maxTokens, refillRate int) *Limiter {
	return &Limiter{
		client: client,
		max:    maxTokens,
		refill: refillRate,
	}
}

// Allow returns true if request is permitted
func (l *Limiter) Allow(ctx context.Context, clientID string) (bool, error) {
	now := time.Now().Unix()

	bucketKey := "bucket:" + clientID
	timestampKey := "bucket_ts:" + clientID

	res, err := l.client.Eval(
		ctx,
		luaScript,
		[]string{bucketKey, timestampKey},
		l.max,
		l.refill,
		now,
	).Int()

	if err != nil {
		return false, err
	}

	return res == 1, nil
}

// Update limiter live (UI control)
func (l *Limiter) Update(bucket, refill int) {
	l.max = bucket
	l.refill = refill
}

// Middleware that enforces rate limit AND logs metrics
func (l *Limiter) Middleware(extract metrics.ClientIDExtractor, m *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {

		clientID := extract(c)
		if clientID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized client"})
			c.Abort()
			return
		}

		ok, err := l.Allow(context.Background(), clientID)

		// Record limiter decision
		m.RecordLimiterDecision(clientID, ok)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "limiter error"})
			c.Abort()
			return
		}
		if !ok {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func (l *Limiter) Max() int     { return l.max }
func (l *Limiter) Refill() int  { return l.refill }
