package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/kuahbanyak/go-crud/pkg/response"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	requests map[string]*bucket
	mu       sync.Mutex
	rate     int           // requests per window
	window   time.Duration // time window
}

type bucket struct {
	tokens   int
	lastSeen time.Time
}

// NewRateLimiter creates a new rate limiter
// rate: number of requests allowed per window
// window: time window duration
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*bucket),
		rate:     rate,
		window:   window,
	}

	// Cleanup old entries every minute
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	b, exists := rl.requests[ip]
	if !exists {
		rl.requests[ip] = &bucket{
			tokens:   rl.rate - 1,
			lastSeen: now,
		}
		return true
	}

	// Reset bucket if window has passed
	if now.Sub(b.lastSeen) > rl.window {
		b.tokens = rl.rate - 1
		b.lastSeen = now
		return true
	}

	// Check if tokens available
	if b.tokens > 0 {
		b.tokens--
		b.lastSeen = now
		return true
	}

	return false
}

// cleanup removes old entries to prevent memory leak
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, b := range rl.requests {
			if now.Sub(b.lastSeen) > rl.window*2 {
				delete(rl.requests, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit middleware limits requests per IP
func RateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract IP from request
			ip := r.RemoteAddr
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = forwarded
			}

			if !limiter.Allow(ip) {
				response.Error(w, http.StatusTooManyRequests, "Rate limit exceeded. Please try again later.", nil)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
