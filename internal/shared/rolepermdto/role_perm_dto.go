package rolepermdto

import "app-server/internal/domain/entity"

type AddNewRoleRequest struct {
	Role    entity.Role `json:"role"`
	PermIDs []uint      `json:"perm_ids"`
}

type ModifyRoleRequest struct {
	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
	PermIDs  []uint `json:"perm_ids"`
}

type AssignRoleToUserRequest struct {
	UserID  uint   `json:"user_id"`
	RoleIDs []uint `json:"role_ids"`
}

type GetRolePermsResponse struct {
	Role  *entity.Role `json:"role"`
	Perms []uint       `json:"perm"`
}

type GroupedResourcesReponse struct {
	ResourceID   uint              `json:"resource_id"`
	ResourceName string            `json:"resource_name"`
	Actions      []ActionWithPerms `json:"actions"`
}

type ActionWithPerms struct {
	ActionID   uint   `json:"action_id"`
	ActionName string `json:"action_name"`
	PermID     uint   `json:"perm_id"`
	PermName   string `json:"perm_name"`
	PermCode   string `json:"perm_code"`
}

type GetPermByUserIdResult struct {
	PermID       int    `json:"perm_id"`
	PermCode     string `json:"perm_code"`
	ResourceCode string `json:"resource_code"`
	ActionCode   string `json:"action_code"`
}
