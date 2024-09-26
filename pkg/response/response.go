package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse định nghĩa phản hồi chuẩn cho API
type APIResponse struct {
	HttpStatus string      `json:"status"`
	Errors     []Errors    `json:"errors,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// Error định nghĩa lỗi trong APIResponse
type Errors struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success trả về phản hồi thành công
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		HttpStatus: http.StatusText(http.StatusOK),
		Data:       data,
	})
}

// Error trả về lỗi tổng quát
func Error(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, APIResponse{
		HttpStatus: http.StatusText(statusCode),
		Errors: []Errors{
			{
				Code:    code,
				Message: message,
			},
		},
	})
}

// ValidationError trả về lỗi validate
func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, APIResponse{
		HttpStatus: http.StatusText(http.StatusBadRequest),
		Errors: []Errors{
			{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		},
	})
}
