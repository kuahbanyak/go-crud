package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/kuahbanyak/go-crud/internal/shared/types"
	"github.com/kuahbanyak/go-crud/pkg/response"
)

// EnhancedRateLimiter provides per-user and per-IP rate limiting
type EnhancedRateLimiter struct {
	buckets map[string]*rateBucket
	mu      sync.RWMutex
	config  RateLimitConfig
}

type RateLimitConfig struct {
	UserRequestsPerWindow int
	IPRequestsPerWindow   int
	Window                time.Duration
	CleanupInterval       time.Duration
}

type rateBucket struct {
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

func NewEnhancedRateLimiter(config RateLimitConfig) *EnhancedRateLimiter {
	rl := &EnhancedRateLimiter{
		buckets: make(map[string]*rateBucket),
		config:  config,
	}
	go rl.cleanup()
	return rl
}

func DefaultRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		UserRequestsPerWindow: 1000,
		IPRequestsPerWindow:   100,
		Window:                time.Hour,
		CleanupInterval:       10 * time.Minute,
	}
}

func StrictRateLimitConfig() RateLimitConfig {
	return RateLimitConfig{
		UserRequestsPerWindow: 5,
		IPRequestsPerWindow:   5,
		Window:                15 * time.Minute,
		CleanupInterval:       5 * time.Minute,
	}
}

func (rl *EnhancedRateLimiter) Allow(key string, isAuthenticated bool) (bool, int, time.Duration) {
	rl.mu.RLock()
	bucket, exists := rl.buckets[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		bucket = &rateBucket{
			tokens:     rl.getMaxTokens(isAuthenticated),
			lastRefill: time.Now(),
		}
		rl.buckets[key] = bucket
		rl.mu.Unlock()
	}

	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastRefill)

	if elapsed >= rl.config.Window {
		bucket.tokens = rl.getMaxTokens(isAuthenticated)
		bucket.lastRefill = now
	}

	if bucket.tokens > 0 {
		bucket.tokens--
		resetTime := bucket.lastRefill.Add(rl.config.Window)
		remaining := time.Until(resetTime)
		return true, bucket.tokens, remaining
	}

	resetTime := bucket.lastRefill.Add(rl.config.Window)
	remaining := time.Until(resetTime)
	return false, 0, remaining
}

func (rl *EnhancedRateLimiter) getMaxTokens(isAuthenticated bool) int {
	if isAuthenticated {
		return rl.config.UserRequestsPerWindow
	}
	return rl.config.IPRequestsPerWindow
}

func (rl *EnhancedRateLimiter) cleanup() {
	ticker := time.NewTicker(rl.config.CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.buckets {
			bucket.mu.Lock()
			if now.Sub(bucket.lastRefill) > rl.config.Window*2 {
				delete(rl.buckets, key)
			}
			bucket.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

func EnhancedRateLimit(limiter *EnhancedRateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var key string
			isAuthenticated := false

			if userID, ok := r.Context().Value("id").(types.MSSQLUUID); ok && userID.String() != "00000000-0000-0000-0000-000000000000" {
				key = fmt.Sprintf("user:%s", userID.String())
				isAuthenticated = true
			} else {
				ip := r.RemoteAddr
				if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
					ip = forwarded
				}
				if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
					ip = realIP
				}
				key = fmt.Sprintf("ip:%s", ip)
			}

			allowed, remaining, resetIn := limiter.Allow(key, isAuthenticated)

			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(resetIn).Unix()))

			if isAuthenticated {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.config.UserRequestsPerWindow))
			} else {
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.config.IPRequestsPerWindow))
			}

			if !allowed {
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(resetIn.Seconds())))
				response.ErrorWithContext(r.Context(), w, http.StatusTooManyRequests,
					fmt.Sprintf("Rate limit exceeded. Try again in %s", resetIn.Round(time.Second)),
					map[string]interface{}{
						"retry_after_seconds": int(resetIn.Seconds()),
						"retry_after":         resetIn.Round(time.Second).String(),
					})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

type EndpointLimiter struct {
	limiters map[string]*EnhancedRateLimiter
	mu       sync.RWMutex
}

func NewEndpointLimiter() *EndpointLimiter {
	return &EndpointLimiter{
		limiters: make(map[string]*EnhancedRateLimiter),
	}
}

func (el *EndpointLimiter) AddEndpoint(pattern string, config RateLimitConfig) {
	el.mu.Lock()
	defer el.mu.Unlock()
	el.limiters[pattern] = NewEnhancedRateLimiter(config)
}

func (el *EndpointLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		el.mu.RLock()
		limiter, exists := el.limiters[r.URL.Path]
		el.mu.RUnlock()

		if !exists {
			next.ServeHTTP(w, r)
			return
		}

		EnhancedRateLimit(limiter)(next).ServeHTTP(w, r)
	})
}
