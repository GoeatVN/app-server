package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse định nghĩa phản hồi chuẩn cho API
type APIResponse struct {
	HttpStatus string      `json:"status"`
	Failed     bool        `json:"isFailed"`
	Errors     []Errors    `json:"errors"`
	Data       interface{} `json:"data"`
	TotalRow   int         `json:"totalRow",omitempty`
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
		Failed:     false,
		Data:       data,
		Errors:     nil,
	})
}

// Error trả về lỗi tổng quát
func Error(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, APIResponse{
		HttpStatus: http.StatusText(statusCode),
		Failed:     true,
		Errors: []Errors{
			{
				Code:    code,
				Message: message,
			},
		},
		Data: nil,
	})
}

// ValidationError trả về lỗi validate
func ValidationError(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, APIResponse{
		HttpStatus: http.StatusText(http.StatusBadRequest),
		Failed:     true,
		Errors: []Errors{
			{
				Code:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		},
		Data: nil,
	})
}
