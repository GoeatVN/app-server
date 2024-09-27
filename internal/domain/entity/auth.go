package entity

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
	RoleName string `json:"role_name"`
	RoleCode string `json:"role_code"`
	BaseEntity
}

func (Role) TableNameRole() string {
	return "roles"
}

type Resource struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ResourceName string `json:"resource_name"`
	ResourceCode string `json:"resource_code"`
	BaseEntity
}

func (Resource) TableName() string {
	return "resources"
}

type Action struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	ActionName string `json:"action_name"`
	ActionCode string `json:"action_code"`
	BaseEntity
}

func (Action) TableName() string {
	return "actions"
}

type Permission struct {
	ID             uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	PermissionName string `json:"permission_name"`
	PermissionCode string `json:"permission_code"`
	ResourceID     uint   `json:"resource_id"`
	ActionID       uint   `json:"action_id"`
	BaseEntity
}

func (Permission) TableName() string {
	return "permissions"
}

type UserRole struct {
	UserID uint `json:"user_id"`
	RoleID uint `json:"role_id"`
	BaseEntity
}

func (UserRole) TableNameUserRole() string {
	return "user_roles"
}

type RolePermission struct {
	RoleID       uint `json:"role_id"`
	PermissionID uint `json:"permission_id"`
	BaseEntity
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
