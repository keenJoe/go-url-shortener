package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter 速率限制器
type RateLimiter struct {
	ips   map[string]*rate.Limiter
	mu    sync.RWMutex
	rate  rate.Limit
	burst int
}

// NewRateLimiter 创建新的速率限制器
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	return &RateLimiter{
		ips:   make(map[string]*rate.Limiter),
		rate:  r,
		burst: b,
	}
}

// GetLimiter 获取IP对应的限制器
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter, exists := rl.ips[ip]
	if !exists {
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.ips[ip] = limiter
	}

	return limiter
}

// RateLimit 限流中间件
func RateLimit(rps int, burst int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate.Limit(rps), burst)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		c.Next()
	}
}
