package metrics

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Tracks request latency; limiter will update blocked/allowed stats separately
func (m *Metrics) Middleware(extract ClientIDExtractor) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		lat := time.Since(start).Milliseconds()

		m.mu.Lock()
		if m.AvgLatencyMs == 0 {
			m.AvgLatencyMs = float64(lat)
		} else {
			m.AvgLatencyMs = (m.AvgLatencyMs + float64(lat)) / 2
		}
		m.mu.Unlock()
	}
}
