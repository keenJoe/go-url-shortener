package middleware

// import (
// 	"sync"

// 	"golang.org/x/time/rate"
// )

// // 限流器映射
// var (
// 	limiterMap = make(map[string]*rate.Limiter)
// 	mu         sync.Mutex
// )

// // 获取限流器
// func getLimiter(ip string) *rate.Limiter {
// 	mu.Lock()
// 	defer mu.Unlock()

// 	limiter, exists := limiterMap[ip]
// 	if !exists {
// 		// 每秒允许10个请求，最多允许30个突发请求
// 		limiter = rate.NewLimiter(10, 30)
// 		limiterMap[ip] = limiter
// 	}

// 	return limiter
// }

// RateLimit 限流中间件
// func RateLimit() gin.HandlerFunc {
// 	return func(
// }
