package api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/limiter"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/metrics"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type WSHandler struct {
	Metrics *metrics.Metrics
	Limiter *limiter.Limiter
	Redis   *redis.Client
}

func NewWebSocketHandler(m *metrics.Metrics, l *limiter.Limiter, rdb *redis.Client) *WSHandler {
	return &WSHandler{Metrics: m, Limiter: l, Redis: rdb}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *WSHandler) Handle(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	ctx := context.Background()

	for {
		time.Sleep(500 * time.Millisecond)

	h.Metrics.RLock()
	snap := h.Metrics.Snapshot()
	h.Metrics.RUnlock()


		// show the *same client* traffic simulator uses
		trafficClient := "dashboard"
		currentTokens, _ := h.Redis.Get(ctx, "bucket:"+trafficClient).Int()

		msg := map[string]interface{}{
			"allowed_total":  snap.AllowedTotal,
			"blocked_total":  snap.BlockedTotal,
			"avg_latency_ms": snap.AvgLatencyMs,
			"clients":        snap.Clients,

			"bucket_size":    h.Limiter.Max(),
			"refill_rate":    h.Limiter.Refill(),
			"current_tokens": currentTokens,
		}

		out, _ := json.Marshal(msg)
		ws.WriteMessage(websocket.TextMessage, out)
	}
}
