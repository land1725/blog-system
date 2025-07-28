package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// Logger 请求日志中间件
// 返回值: Gin处理函数
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		log.Printf("Method: %s | Path: %s | Status: %d | Duration: %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			duration,
		)
	}
}
