package middleware

import (
	"flexgate/internal/limiter"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RateLimitMiddleware 路由限流中间件
// 每个路由在初始化时会分配一个专有的 bucket 实例
func RateLimitMiddleware(bucket *limiter.TokenBucket) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 尝试从桶里拿令牌
		if !bucket.Allow() {
			// 拿不到令牌，说明触发限流
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too Many Requests",
				"message": "当前接口访问过于频繁，请稍后再试",
				"code":    429,
			})
			// 关键：调用 Abort() 停止执行后续的 Proxy 逻辑
			c.Abort()
			return
		}

		// 拿到令牌，继续执行下一个中间件或 Proxy 逻辑
		c.Next()
	}
}
