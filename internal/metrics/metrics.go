package metrics

import (
	"sync"
)

type ClientStats struct {
	Allowed int `json:"allowed"`
	Blocked int `json:"blocked"`
}

// This struct is NOT returned directly because it contains a mutex.
type Metrics struct {
    mu sync.RWMutex
    AllowedTotal int
    BlockedTotal int
    AvgLatencyMs float64
    Clients map[string]*ClientStats

    BucketSize  int
    RefillRate  int
}

// Safe read access from outside package
func (m *Metrics) RLock()   { m.mu.RLock() }
func (m *Metrics) RUnlock() { m.mu.RUnlock() }
func (m *Metrics) Lock()    { m.mu.Lock() }
func (m *Metrics) Unlock()  { m.mu.Unlock() }


// Snapshot struct is safe for returning via WebSocket
type Snapshot struct {
	AllowedTotal  int                    `json:"allowed_total"`
	BlockedTotal  int                    `json:"blocked_total"`
	AvgLatencyMs  float64                `json:"avg_latency_ms"`
	BucketSize    int                    `json:"bucket_size"`
	RefillRate    int                    `json:"refill_rate"`
	Clients       map[string]*ClientStats `json:"clients"`
}

// Constructor
func NewMetrics() *Metrics {
	return &Metrics{
		Clients: make(map[string]*ClientStats),
	}
}

// Called by Limiter middleware
func (m *Metrics) RecordLimiterDecision(clientID string, allowed bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.Clients[clientID]; !exists {
		m.Clients[clientID] = &ClientStats{}
	}

	if allowed {
		m.AllowedTotal++
		m.Clients[clientID].Allowed++
	} else {
		m.BlockedTotal++
		m.Clients[clientID].Blocked++
	}
}

// Called by websocket
func (m *Metrics) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	out := Snapshot{
		AllowedTotal: m.AllowedTotal,
		BlockedTotal: m.BlockedTotal,
		AvgLatencyMs: m.AvgLatencyMs,
		BucketSize:   m.BucketSize,
		RefillRate:   m.RefillRate,
		Clients:      make(map[string]*ClientStats),
	}

	for k, v := range m.Clients {
		out.Clients[k] = &ClientStats{Allowed: v.Allowed, Blocked: v.Blocked}
	}

	return out
}

// Shared extractor type (used by limiter + metrics middleware)
