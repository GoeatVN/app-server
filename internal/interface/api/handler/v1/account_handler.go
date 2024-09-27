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

	var loginDto login.LoginDTO

	if err := c.ShouldBindJSON(&loginDto); err != nil {
		response.ValidationError(c, err)
		return
	}

	token, err := h.accountService.Login(loginDto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Set("response_data", token)
}
