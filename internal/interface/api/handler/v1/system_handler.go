package v1

import (
	"app-server/internal/shared/systemdto"
	"app-server/internal/usecase/system"
	"app-server/pkg/response"
	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	service system.SystemServiceInterface
}

func NewSystemHandler(service system.SystemServiceInterface) *SystemHandler {
	return &SystemHandler{service: service}
}

func (h *SystemHandler) LoadComboboxDataHandler(c *gin.Context) {
	var request systemdto.ComboboxRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		// Gọi trực tiếp response.ValidationError
		response.ValidationError(c, err)
		return
	}

	// Call the service method
	response, err := h.service.LoadComboboxData(request)
	if err != nil {
		c.Error(err)
		return
	}

	// Return the response
	c.Set("response_data", response)
}
