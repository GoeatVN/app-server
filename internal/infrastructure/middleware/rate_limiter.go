package middleware

import (
	"net/http"
	"time"

	"app-server/pkg/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	rateLimiters = make(map[string]*rate.Limiter)
	rateLimit    = 1000        // Số lượng yêu cầu tối đa
	burstLimit   = 10          // Số lượng yêu cầu bùng phát tối đa
	rateDuration = time.Minute // Thời gian giới hạn
)

// getClientIP lấy địa chỉ IP của client
func getClientIP(c *gin.Context) string {
	clientIP := c.ClientIP()
	if clientIP == "" {
		clientIP = c.Request.RemoteAddr
	}
	return clientIP
}

// RateLimiterMiddleware giới hạn số lượng yêu cầu từ một địa chỉ IP
func RateLimiterMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := getClientIP(c)

		// Kiểm tra xem IP có trong map rate limiters không
		if limiter, exists := rateLimiters[clientIP]; exists {
			// Nếu đã vượt quá giới hạn, trả về lỗi 429
			if !limiter.Allow() {
				response.Error(c, http.StatusTooManyRequests, "TOO_MANY_REQUESTS", "Rate limit exceeded. Please try again later.")
				c.Abort()
				return
			}
		} else {
			// Tạo một limiter mới cho IP
			limiter := rate.NewLimiter(rate.Every(rateDuration/time.Duration(rateLimit)), burstLimit)
			rateLimiters[clientIP] = limiter
		}

		c.Next() // Tiếp tục xử lý request nếu chưa vượt giới hạn
	}
}
