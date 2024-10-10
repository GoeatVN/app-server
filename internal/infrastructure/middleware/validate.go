package middleware

import (
	"app-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// ValidateMiddleware kiểm tra tính hợp lệ của dữ liệu đầu vào
func ValidateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name  string `json:"name" binding:"required"`
			Email string `json:"email" binding:"required,email"`
		}

		// Kiểm tra xem dữ liệu nhập vào có hợp lệ không
		if err := c.ShouldBindJSON(&input); err != nil {
			response.ValidationError(c, err)
			c.Abort()
			return
		}

		c.Next()
	}
}
