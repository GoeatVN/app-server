package middleware

import (
	"app-server/internal/usecase/auth"
	"app-server/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AuthorizationMiddleware kiểm tra quyền truy cập dựa trên vai trò người dùng
func AuthorizationMiddleware(permission string, authService auth.AuthServiceInterface) gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		_, err := authService.GetClaims(authHeader)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
			c.Abort()
			return
		}

		// Kiểm tra nếu vai trò của người dùng không khớp với yêu cầu
		//if userRole != permission {
		//	response.Error(c, http.StatusForbidden, "FORBIDDEN", "You do not have permission to access this resource")
		//	c.Abort() // Dừng xử lý request
		//	return
		//}

		// Nếu vai trò hợp lệ, tiếp tục xử lý request
		c.Next()
	}
}
