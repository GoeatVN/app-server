package v1

import (
	"app-server/internal/shared/login"
	"app-server/internal/usecase/account"
	"app-server/pkg/response"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService account.ServiceInterface
}

func NewAccountHandler(accountService account.ServiceInterface) *AccountHandler {
	return &AccountHandler{accountService: accountService}
}

// Login handles user login
func (h *AccountHandler) Login(c *gin.Context) {

	var loginRequest login.LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		response.ValidationError(c, err)
		return
	}

	loginResponse, err := h.accountService.Login(loginRequest)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("response_data", loginResponse)
}
