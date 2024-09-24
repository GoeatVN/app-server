package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware logs the details of each request
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy thời gian bắt đầu
		startTime := time.Now()

		// Xử lý yêu cầu
		c.Next()

		// Lấy thời gian kết thúc
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)

		// Lấy thông tin yêu cầu
		method := c.Request.Method
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		path := c.Request.URL.Path

		// Ghi log
		log.Printf("| %3d | %13v | %15s | %s  %s",
			statusCode,
			latencyTime,
			clientIP,
			method,
			path,
		)
	}
}
