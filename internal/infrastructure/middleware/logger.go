package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware ghi lại log của các yêu cầu HTTP
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		// Ghi log sau khi xử lý xong
		log.Printf("Method: %s | Path: %s | Status: %d | Latency: %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(startTime),
		)
	}
}
