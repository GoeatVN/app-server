package rolepermdto

import "app-server/internal/domain/entity"

type AddNewRoleRequest struct {
	GetRolePermsResponse
}

type ModifyRoleRequest struct {
	GetRolePermsResponse
}

type AssignRoleToUserRequest struct {
	UserID  uint   `json:"user_id"`
	RoleIDs []uint `json:"role_ids"`
}

type GetRolePermsResponse struct {
	Role  *entity.Role `json:"role"`
	Perms []uint       `json:"perms"`
}

type GroupedResourcesReponse struct {
	ResourceID   uint              `json:"resourceId"`
	ResourceName string            `json:"resourceName"`
	Actions      []ActionWithPerms `json:"actions"`
}

type ActionWithPerms struct {
	ActionID   uint   `json:"actionId"`
	ActionName string `json:"actionName"`
	PermID     uint   `json:"permId"`
	PermName   string `json:"permName"`
	PermCode   string `json:"permCode"`
}

type GetPermByUserIdResult struct {
	PermID       int    `json:"permId"`
	PermCode     string `json:"permCode"`
	ResourceCode string `json:"resourceCode"`
	ActionCode   string `json:"actionCode"`
}
