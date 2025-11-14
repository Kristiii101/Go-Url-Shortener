package observability

import (
	"sync"
	"time"
)

// Metrics tracks simple application metrics
type Metrics struct {
	mu             sync.RWMutex
	LinksCreated   int64
	LinksClicked   int64
	RequestsServed int64
	ErrorsCount    int64
	StartTime      time.Time
	LastResetTime  time.Time
}

// NewMetrics creates a new metrics instance
func NewMetrics() *Metrics {
	now := time.Now()
	return &Metrics{
		StartTime:     now,
		LastResetTime: now,
	}
}

// IncrementLinksCreated increments the links created counter
func (m *Metrics) IncrementLinksCreated() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LinksCreated++
}

// IncrementLinksClicked increments the links clicked counter
func (m *Metrics) IncrementLinksClicked() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.LinksClicked++
}

// IncrementRequests increments the requests counter
func (m *Metrics) IncrementRequests() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RequestsServed++
}

// IncrementErrors increments the errors counter
func (m *Metrics) IncrementErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ErrorsCount++
}

// GetStats returns a copy of current stats
func (m *Metrics) GetStats() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	uptime := time.Since(m.StartTime)

	return map[string]interface{}{
		"links_created":   m.LinksCreated,
		"links_clicked":   m.LinksClicked,
		"requests_served": m.RequestsServed,
		"errors_count":    m.ErrorsCount,
		"uptime_seconds":  uptime.Seconds(),
		"start_time":      m.StartTime,
	}
}

// Reset resets all counters
func (m *Metrics) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.LinksCreated = 0
	m.LinksClicked = 0
	m.RequestsServed = 0
	m.ErrorsCount = 0
	m.LastResetTime = time.Now()
}
