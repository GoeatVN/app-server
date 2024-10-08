package middleware

import (
	"app-server/internal/shared/constants"
	"app-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// ResponseHandlerMiddleware đảm bảo mọi phản hồi đều có định dạng chuẩn
func ResponseHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// // Kiểm tra xem có lỗi xảy ra trong quá trình xử lý không
		// if len(c.Errors) > 0 {
		// 	// Xử lý tất cả các lỗi trả về
		// 	response.Error(c, c.Writer.Status(), "INTERNAL_ERROR", c.Errors.String())
		// 	return
		// }

		// Nếu không có lỗi, chuẩn hóa phản hồi thành công (nếu có)

		if data, exists := c.Get(constants.RESPONSE_DATA_KEY); exists {
			response.Success(c, data)
		}
	}
}
