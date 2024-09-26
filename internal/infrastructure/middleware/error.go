package middleware

import (
	"net/http"

	"app-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler là middleware xử lý lỗi
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Xử lý request tiếp theo trong chain
		c.Next()

		// Kiểm tra nếu có lỗi
		if len(c.Errors) > 0 {
			// Sử dụng response.Error để định dạng lỗi và trả về phản hồi
			response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", c.Errors.String())
			c.Abort() // Ngừng quá trình xử lý tiếp theo
			return
		}
	}
}
