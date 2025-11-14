package rate

import (
	"sync"
	"time"
)

// Limiter provides rate limiting functionality
type Limiter struct {
	mu       sync.RWMutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

type visitor struct {
	count     int
	firstSeen time.Time
}

// NewLimiter creates a new rate limiter
func NewLimiter(requestsPerWindow int, window time.Duration) *Limiter {
	l := &Limiter{
		visitors: make(map[string]*visitor),
		limit:    requestsPerWindow,
		window:   window,
	}

	// Clean up old entries periodically
	go l.cleanup()

	return l
}

// Allow checks if a request from the given key is allowed
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	v, exists := l.visitors[key]

	if !exists {
		l.visitors[key] = &visitor{
			count:     1,
			firstSeen: now,
		}
		return true
	}

	// Reset if window has passed
	if now.Sub(v.firstSeen) > l.window {
		v.count = 1
		v.firstSeen = now
		return true
	}

	// Check limit
	if v.count >= l.limit {
		return false
	}

	v.count++
	return true
}

// cleanup removes old entries
func (l *Limiter) cleanup() {
	ticker := time.NewTicker(l.window)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for key, v := range l.visitors {
			if now.Sub(v.firstSeen) > l.window*2 {
				delete(l.visitors, key)
			}
		}
		l.mu.Unlock()
	}
}
