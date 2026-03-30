package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 是一个自定义的日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 请求开始时间
		startTime := time.Now()

		// 2. 处理请求（执行后续的代理转发逻辑）
		c.Next()

		// 3. 请求结束时间
		endTime := time.Now()

		// 4. 计算执行耗时
		latencyTime := endTime.Sub(startTime)

		// 5. 获取请求信息
		reqMethod := c.Request.Method
		reqUri := c.Request.RequestURI
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		// 6. 打印结构化日志
		// 格式：[FlexGate] 200 | 1.25ms | 127.0.0.1 | GET "/api/users"
		log.Printf("[FlexGate] %3d | %13v | %15s | %s %s",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
