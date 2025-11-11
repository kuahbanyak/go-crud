package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/kuahbanyak/go-crud/pkg/response"
)

type RateLimiter struct {
	requests map[string]*bucket
	mu       sync.Mutex
	rate     int
	window   time.Duration
}

type bucket struct {
	tokens   int
	lastSeen time.Time
}

func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string]*bucket),
		rate:     rate,
		window:   window,
	}

	go rl.cleanup()

	return rl
}

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

	if now.Sub(b.lastSeen) > rl.window {
		b.tokens = rl.rate - 1
		b.lastSeen = now
		return true
	}

	if b.tokens > 0 {
		b.tokens--
		b.lastSeen = now
		return true
	}

	return false
}

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

func RateLimit(limiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
