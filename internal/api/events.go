package api

import (
	"context"
	"net"
	"net/http"

	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/models"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	store *storage.EventStore
}

func NewEventHandler(store *storage.EventStore) *EventHandler {
	return &EventHandler{store: store}
}

func clientIP(c *gin.Context) string {
	ip, _, _ := net.SplitHostPort(c.Request.RemoteAddr)
	if ip == "" {
		return c.ClientIP()
	}
	return ip
}

func (h *EventHandler) Ingest(c *gin.Context) {
	var body struct {
		EventID string                 `json:"event_id"`
		Source  string                 `json:"source"`
		Payload map[string]interface{} `json:"payload"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if body.EventID == "" || body.Source == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event_id and source required"})
		return
	}

	id, err := uuid.Parse(body.EventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event_id"})
		return
	}

	e := models.Event{
		ID:       id,
		Source:   body.Source,
		Payload:  body.Payload,
		ClientIP: clientIP(c),
	}

	if err := h.store.InsertEvent(context.Background(), e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db insert failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "accepted"})
}
