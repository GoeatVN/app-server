package middleware

import (
	"app-server/internal/usecase/auth"
	"net/http"
	"strings"

	"app-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware kiểm tra Bearer token từ header HTTP
func AuthenticationMiddleware(authService auth.AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Kiểm tra header Authorization có chứa Bearer token không
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid token")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Kiểm tra token bằng cách sử dụng authService
		_, err := authService.VerifyToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
			c.Abort()
			return
		}
		c.Next()
	}
}
