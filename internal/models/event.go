package models

import "github.com/google/uuid"

type Event struct {
	ID        uuid.UUID              `json:"event_id"`
	Source    string                 `json:"source"`
	Payload   map[string]interface{} `json:"payload"`
	ClientIP  string                 `json:"-"`
}
