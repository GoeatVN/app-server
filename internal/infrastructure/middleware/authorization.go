package middleware

import (
	"app-server/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthorizationMiddleware kiểm tra quyền truy cập dựa trên vai trò người dùng
func AuthorizationMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Giả định rằng vai trò được gửi qua header Authorization hoặc Bearer token
		userRole := c.GetHeader("Role") // Có thể thay bằng cách kiểm tra từ token nếu cần

		// Kiểm tra nếu vai trò của người dùng không khớp với yêu cầu
		if userRole != requiredRole {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "You do not have permission to access this resource")
			c.Abort() // Dừng xử lý request
			return
		}

		// Nếu vai trò hợp lệ, tiếp tục xử lý request
		c.Next()
	}
}
