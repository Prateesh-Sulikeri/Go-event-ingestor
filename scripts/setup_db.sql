-- Event Storage Schema

CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    source VARCHAR(128) NOT NULL,
    payload JSONB NOT NULL,
    client_ip INET,
    received_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(id)
);

-- Indexing for analytical queries and dashboard
CREATE INDEX IF NOT EXISTS idx_events_received_at ON events (received_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_source ON events (source);

