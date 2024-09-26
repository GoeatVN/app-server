package v1

import (
	"app-server/internal/domain/entity"
	"app-server/internal/usecase/user"
	"app-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service user.ServiceInterface
}

func NewUserHandler(service user.ServiceInterface) *UserHandler {
	return &UserHandler{service: service}
}

// Lấy danh sách người dùng
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		// Ghi lại lỗi vào context và để ErrorHandler xử lý
		c.Error(err)
		return
	}
	// Đặt dữ liệu phản hồi vào context để ResponseHandlerMiddleware xử lý
	c.Set("response_data", users)
}

// Tạo người dùng mới
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		// Gọi trực tiếp response.ValidationError
		response.ValidationError(c, err)
		return
	}
	err := h.service.CreateUser(&user)
	if err != nil {
		// Ghi lại lỗi vào context để ErrorHandler xử lý
		c.Error(err)
		return
	}
	// Đặt dữ liệu phản hồi vào context để ResponseHandlerMiddleware xử lý
	c.Set("response_data", user)
}
