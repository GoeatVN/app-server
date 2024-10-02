package v1

import (
	"app-server/internal/shared/rolepermdto"
	"app-server/internal/usecase/rolepermission"
	"app-server/pkg/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

type RolePermHandler struct {
	rolePermService rolepermission.RolePermServiceInterface
}

func NewRolePermHandler(rolePermService rolepermission.RolePermServiceInterface) *RolePermHandler {
	return &RolePermHandler{rolePermService: rolePermService}
}

// AddNewRole handles adding a new role
func (h *RolePermHandler) AddNewRole(c *gin.Context) {
	var addNewRoleRequest rolepermdto.AddNewRoleRequest
	if err := c.ShouldBindJSON(&addNewRoleRequest); err != nil {
		response.ValidationError(c, err)
		return
	}
	if err := h.rolePermService.AddNewRole(addNewRoleRequest); err != nil {
		c.Error(err)
		return
	}
	c.Set("response_data", "Role added successfully")
}

// ModifyRole handles modifying a role
func (h *RolePermHandler) ModifyRole(c *gin.Context) {
	var modifyRoleRequest rolepermdto.ModifyRoleRequest
	if err := c.ShouldBindJSON(&modifyRoleRequest); err != nil {
		response.ValidationError(c, err)
		return
	}
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(err)
		return
	}
	if err := h.rolePermService.ModifyRole(uint(id), modifyRoleRequest); err != nil {
		c.Error(err)
		return
	}
	c.Set("response_data", "Role modified successfully")
}

// AssignRoleToUser handles assigning roles to a user
func (h *RolePermHandler) AssignRoleToUser(c *gin.Context) {
	var assignRoleToUserRequest rolepermdto.AssignRoleToUserRequest
	if err := c.ShouldBindJSON(&assignRoleToUserRequest); err != nil {
		response.ValidationError(c, err)
		return
	}
	if err := h.rolePermService.AssignRoleToUser(assignRoleToUserRequest); err != nil {
		c.Error(err)
		return
	}
	c.Set("response_data", "Role assigned successfully")
}

func (h *RolePermHandler) GetAllRolePerms(c *gin.Context) {
	perms, err := h.rolePermService.GetRolePerms(0)
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("response_data", perms)
}

func (h *RolePermHandler) GetRolePermsById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(err)
		return
	}
	perms, err := h.rolePermService.GetRolePerms(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("response_data", perms)
}

func (h *RolePermHandler) GetRoleGroupByResource(c *gin.Context) {
	perms, err := h.rolePermService.GetRoleGroupByResource()
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("response_data", perms)
}

func (h *RolePermHandler) GetPermsByUserID(c *gin.Context) {

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Error(err)
		return
	}

	perms, err := h.rolePermService.GetPermsByUserID(uint(id))
	if err != nil {
		c.Error(err)
		return
	}
	c.Set("response_data", perms)
}
