package middleware

import (
	"net/http"
	"sync"
	"time"

	"go-crud/pkg/response"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     time.Duration
	capacity int
}

type visitor struct {
	lastSeen time.Time
	tokens   int
}

func NewRateLimiter(rate time.Duration, capacity int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		capacity: capacity,
	}

	go rl.cleanupVisitors()

	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{
			lastSeen: time.Now(),
			tokens:   rl.capacity - 1,
		}
		return true
	}

	now := time.Now()
	elapsed := now.Sub(v.lastSeen)
	tokensToAdd := int(elapsed / rl.rate)

	if tokensToAdd > 0 {
		v.tokens += tokensToAdd
		if v.tokens > rl.capacity {
			v.tokens = rl.capacity
		}
		v.lastSeen = now
	}

	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > time.Hour {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func RateLimit(requestsPerMinute int) gin.HandlerFunc {
	limiter := NewRateLimiter(time.Minute/time.Duration(requestsPerMinute), requestsPerMinute)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded", "Too many requests from this IP")
			c.Abort()
			return
		}
		c.Next()
	}
}
