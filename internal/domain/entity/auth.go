package entity

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	BaseEntity
}

func (User) TableName() string {
	return "users"
}

type Role struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleName string `json:"roleName"`
	RoleCode string `json:"roleCode"`
	BaseEntity
}

func (Role) TableNameRole() string {
	return "roles"
}

type Resource struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceName string `json:"resourceName"`
	ResourceCode string `json:"resourceCode"`
	BaseEntity
}

func (Resource) TableName() string {
	return "resources"
}

type Action struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ActionName string `json:"actionName"`
	ActionCode string `json:"actionCode"`
	BaseEntity
}

func (Action) TableName() string {
	return "actions"
}

type Permission struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PermissionName string `json:"permissionName"`
	PermissionCode string `json:"permissionCode"`
	ResourceID     uint   `json:"resourceId"`
	ActionID       uint   `json:"actionId"`
	BaseEntity
}

func (Permission) TableName() string {
	return "permissions"
}

type UserRole struct {
	UserID uint `json:"userId"`
	RoleID uint `json:"roleId"`
	BaseEntity
}

func (UserRole) TableNameUserRole() string {
	return "user_roles"
}

type RolePermission struct {
	RoleID       uint `json:"roleId"`
	PermissionID uint `json:"permissionId"`
	BaseEntity
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// Struct Claims chứa thông tin của người dùng và các quyền
type AuthClaims struct {
	UserID   uint   `json:"userId"`
	RoleIDs  []uint `json:"roleIds"`
	Username string `json:"username"`
	jwt.StandardClaims
}
