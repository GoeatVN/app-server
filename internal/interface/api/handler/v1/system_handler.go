package v1

import (
	"app-server/internal/shared/constants"
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
	var request []systemdto.ComboboxRequestItem

	if err := c.ShouldBindJSON(&request); err != nil {
		// Gọi trực tiếp combodata.ValidationError
		response.ValidationError(c, err)
		return
	}

	// Call the service method
	combodata, err := h.service.LoadComboboxData(request)
	if err != nil {
		c.Error(err)
		return
	}

	// Return the combodata
	c.Set(constants.RESPONSE_DATA_KEY, combodata)
}
