package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/config"
	"github.com/Prateesh-Sulikeri/Go-event-ingestor/internal/models"
	_ "github.com/lib/pq"
)

type EventStore struct {
	db *sql.DB
}

func NewEventStore(cfg config.Config) (*EventStore, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB_HOST, cfg.DB_PORT, cfg.DB_USER, cfg.DB_PASSWORD, cfg.DB_NAME,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &EventStore{db: db}, nil
}

func (s *EventStore) InsertEvent(ctx context.Context, e models.Event) error {
	payloadBytes, err := json.Marshal(e.Payload)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(
		ctx,
		`INSERT INTO events (id, source, payload, client_ip)
		 VALUES ($1, $2, $3::jsonb, $4)
		 ON CONFLICT (id) DO NOTHING`,
		e.ID, e.Source, string(payloadBytes), e.ClientIP,
	)

	return err
}
