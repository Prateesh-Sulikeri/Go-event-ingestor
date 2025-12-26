package limiter

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClientIDExtractor func(*gin.Context) string

// Middleware returns gin.HandlerFunc with custom client ID extraction.
func (l *Limiter) Middleware(extract ClientIDExtractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := extract(c)
		if clientID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized client"})
			c.Abort()
			return
		}

		ok, err := l.Allow(context.Background(), clientID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "limiter error"})
			c.Abort()
			return
		}

		if !ok {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
