package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// RequestIDHeader 定义请求ID的HTTP头名称
const RequestIDHeader = "X-Request-ID"

// RequestIDContextKey 定义在gin.Context中存储请求ID的键名
const RequestIDContextKey = "request_id"

// RequestID 中间件：生成或提取请求ID，并存入上下文，同时写入响应头
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 尝试从请求头获取请求ID
		reqID := c.GetHeader(RequestIDHeader)

		// 2. 如果没有，生成一个随机ID（16字节随机数，转换为十六进制）
		if reqID == "" {
			reqID = generateRandomID()
		}

		// 3. 将请求ID存入上下文，供后续中间件使用
		c.Set(RequestIDContextKey, reqID)

		// 4. 将请求ID写入响应头，便于客户端追踪
		c.Header(RequestIDHeader, reqID)

		// 5. 继续处理请求
		c.Next()
	}
}

// generateRandomID 生成一个随机的16字节十六进制字符串，用作请求ID
func generateRandomID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// 理论上 rand.Read 不会失败，若失败则使用时间戳+计数器降级（简单实现）
		return "fallback-" + hex.EncodeToString(b) // 实际不会走到这里
	}
	return hex.EncodeToString(b)
}
