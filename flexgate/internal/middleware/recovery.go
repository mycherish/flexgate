package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 返回一个 Gin 中间件，用于捕获处理过程中的 panic，
// 记录错误日志和堆栈信息，并返回 500 状态码的 JSON 响应。
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 记录错误信息
				log.Printf("[Recovery] panic recovered: %v\n%s", err, debug.Stack())

				// 终止后续处理并返回统一错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    500,
					"error":   "Internal Server Error",
					"message": fmt.Sprintf("%v", err),
				})
			}
		}()
		c.Next()
	}
}
