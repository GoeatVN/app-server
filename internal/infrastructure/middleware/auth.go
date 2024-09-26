package middleware

import (
	"net/http"
	"strings"

	"app-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware kiểm tra Bearer token từ header HTTP
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Kiểm tra header Authorization có chứa Bearer token không
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid token")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Kiểm tra token (bạn có thể thay bằng logic thực sự)
		if !isValidToken(token) {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
			c.Abort()
			return
		}

		c.Next()
	}
}

// Giả lập kiểm tra token
func isValidToken(token string) bool {
	// Thay thế bằng logic kiểm tra token thực sự
	return token == "valid_token"
}
