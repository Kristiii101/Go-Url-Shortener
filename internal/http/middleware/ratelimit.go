package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// simple token-bucket per IP (in-memory)
type bucket struct {
	mu         sync.Mutex
	tokens     float64
	capacity   float64
	fillRate   float64 // tokens per second
	lastRefill time.Time
}

func newBucket(capacity int, window time.Duration) *bucket {
	return &bucket{
		tokens:     float64(capacity),
		capacity:   float64(capacity),
		fillRate:   float64(capacity) / window.Seconds(),
		lastRefill: time.Now(),
	}
}

func (b *bucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(b.lastRefill).Seconds()
	// Refill
	b.tokens = minF(b.capacity, b.tokens+elapsed*b.fillRate)
	b.lastRefill = now
	if b.tokens >= 1.0 {
		b.tokens -= 1.0
		return true
	}
	return false
}

func minF(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

type limiter struct {
	mu      sync.Mutex
	buckets map[string]*bucket
	cap     int
	win     time.Duration
}

func newLimiter(cap int, window time.Duration) *limiter {
	return &limiter{
		buckets: make(map[string]*bucket),
		cap:     cap,
		win:     window,
	}
}

func (l *limiter) get(ip string) *bucket {
	l.mu.Lock()
	defer l.mu.Unlock()
	if b, ok := l.buckets[ip]; ok {
		return b
	}
	b := newBucket(l.cap, l.win)
	l.buckets[ip] = b
	return b
}

func RateLimitPerIP(cap int, window time.Duration) func(http.Handler) http.Handler {
	l := newLimiter(cap, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			if ip == "" {
				ip = "unknown"
			}
			if !l.get(ip).allow() {
				w.Header().Set("Retry-After", window.String())
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	// For local dev, RemoteAddr is fine (strip port)
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}
