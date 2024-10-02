package rolepermission

import (
	"app-server/internal/domain/entity"
	"app-server/internal/persistence/repository"
	"app-server/internal/persistence/repository/postgres"
	"app-server/internal/shared/rolepermdto"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

type rolePermService struct {
	userRepo       *postgres.UserRepository
	userRoleRepo   *repository.GenericBaseRepository[entity.UserRole]
	roleRepo       *repository.GenericBaseRepository[entity.Role]
	rolePermission *repository.GenericBaseRepository[entity.RolePermission]
	permRepo       *repository.GenericBaseRepository[entity.Permission]
	resourceRepo   *repository.GenericBaseRepository[entity.Resource]
	actionRepo     *repository.GenericBaseRepository[entity.Action]
	db             *gorm.DB
	pool           *pgxpool.Pool
}

type RolePermServiceInterface interface {
	AddNewRole(request rolepermdto.AddNewRoleRequest) error
	ModifyRole(roleId uint, request rolepermdto.ModifyRoleRequest) error
	AssignRoleToUser(request rolepermdto.AssignRoleToUserRequest) error
	GetRolePerms(id uint) ([]rolepermdto.GetRolePermsResponse, error)
	GetRoleGroupByResource() ([]rolepermdto.GroupedResourcesReponse, error)
	GetPermsByUserID(userID uint) ([]rolepermdto.GetPermByUserIdResult, error)
}

func NewRolePermService(userRepo *postgres.UserRepository,
	userRoleRepo *repository.GenericBaseRepository[entity.UserRole],
	roleRepo *repository.GenericBaseRepository[entity.Role],
	rolePermission *repository.GenericBaseRepository[entity.RolePermission],
	permRepo *repository.GenericBaseRepository[entity.Permission],
	resourceRepo *repository.GenericBaseRepository[entity.Resource],
	actionRepo *repository.GenericBaseRepository[entity.Action],
	db *gorm.DB,
	pool *pgxpool.Pool,
) RolePermServiceInterface {
	return &rolePermService{userRepo: userRepo,
		userRoleRepo: userRoleRepo, roleRepo: roleRepo,
		rolePermission: rolePermission,
		permRepo:       permRepo, resourceRepo: resourceRepo,
		actionRepo: actionRepo, db: db, pool: pool}
}

func (s *rolePermService) AddNewRole(request rolepermdto.AddNewRoleRequest) error {
	// Create the role
	if err := s.roleRepo.Create(&request.Role); err != nil {
		return err
	}

	// Assign permissions to the role
	for _, permID := range request.PermIDs {
		rolePerm := entity.RolePermission{RoleID: request.Role.ID, PermissionID: permID}
		if err := s.rolePermission.Create(&rolePerm); err != nil {
			return err
		}
	}
	return nil
}

func (s *rolePermService) ModifyRole(roleId uint, request rolepermdto.ModifyRoleRequest) error {
	role, err := s.roleRepo.FindByID(roleId)
	if err != nil {
		return err
	}

	role.RoleName = request.RoleName

	if err := s.roleRepo.Update(role); err != nil {
		return err
	}

	// Retrieve existing role permissions
	var existingRolePerms []entity.RolePermission
	if err := s.rolePermission.Where("role_id = ?", roleId).Find(&existingRolePerms).Error; err != nil {
		return err
	}

	existingPermIDs := make(map[uint]bool)
	for _, rp := range existingRolePerms {
		existingPermIDs[rp.PermissionID] = true
	}

	// Insert new permissions
	for _, permID := range request.PermIDs {
		if !existingPermIDs[permID] {
			rolePerm := entity.RolePermission{RoleID: roleId, PermissionID: permID}
			if err := s.rolePermission.Create(&rolePerm); err != nil {
				return err
			}
		}
		delete(existingPermIDs, permID)
	}

	// Delete permissions that are no longer needed
	for permID := range existingPermIDs {
		if err := s.rolePermission.Where("role_id = ? AND permission_id = ?", roleId, permID).Delete(&entity.RolePermission{}).Error; err != nil {
			return err
		}
	}

	return nil
}

func (s *rolePermService) AssignRoleToUser(request rolepermdto.AssignRoleToUserRequest) error {
	// Retrieve existing user roles
	var existingUserRoles []entity.UserRole
	if err := s.userRoleRepo.Where("user_id = ?", request.UserID).Find(&existingUserRoles).Error; err != nil {
		return err
	}

	existingRoleIDs := make(map[uint]bool)
	for _, ur := range existingUserRoles {
		existingRoleIDs[ur.RoleID] = true
	}

	// Insert new roles
	for _, roleID := range request.RoleIDs {
		if !existingRoleIDs[roleID] {
			userRole := entity.UserRole{UserID: request.UserID, RoleID: roleID}
			if err := s.userRoleRepo.Create(&userRole); err != nil {
				return err
			}
		}
		delete(existingRoleIDs, roleID)
	}

	// Delete roles that are no longer needed
	for roleID := range existingRoleIDs {
		if err := s.userRoleRepo.Where("user_id = ? AND role_id = ?", request.UserID, roleID).Delete(&entity.UserRole{}).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetRoles handles getting all roles or a specific role by ID if provided
func (s *rolePermService) GetRolePerms(id uint) ([]rolepermdto.GetRolePermsResponse, error) {
	var rows []struct {
		RoleID   uint
		RoleName string
		PermID   uint
	}

	// Truy vấn SQL để lấy thông tin role và permission
	query := `
		SELECT 
			r.id AS role_id, r.role_name, rp.permission_id AS perm_id
		FROM roles r
		LEFT JOIN role_permissions rp ON rp.role_id = r.id
		WHERE (? = 0 OR r.id = ?)
	`

	// Thực thi truy vấn
	if err := s.db.Raw(query, id, id).Scan(&rows).Error; err != nil {
		return nil, err
	}

	// Tạo map để gom nhóm permissions theo role
	roleMap := make(map[uint]*rolepermdto.GetRolePermsResponse)
	for _, row := range rows {
		// Nếu role chưa tồn tại trong map, tạo mới
		if _, ok := roleMap[row.RoleID]; !ok {
			roleMap[row.RoleID] = &rolepermdto.GetRolePermsResponse{
				Role: &entity.Role{
					ID:       row.RoleID,
					RoleName: row.RoleName,
				},
				Perms: []uint{},
			}
		}

		// Thêm permission ID vào danh sách
		roleMap[row.RoleID].Perms = append(roleMap[row.RoleID].Perms, row.PermID)
	}

	// Chuyển từ map sang slice để trả về
	var result []rolepermdto.GetRolePermsResponse
	for _, role := range roleMap {
		result = append(result, *role)
	}

	return result, nil
}

func (s *rolePermService) GetRoleGroupByResource() ([]rolepermdto.GroupedResourcesReponse, error) {
	var rows []struct {
		ResourceID   uint
		ResourceName string
		ActionID     uint
		ActionName   string
		PermID       uint
		PermName     string
		PermCode     string
	}

	// SQL query to fetch permissions, resources, and actions
	query := `
		SELECT
			p.id AS perm_id, p.permission_name AS perm_name, p.permission_code AS perm_code,
			r.id AS resource_id, r.resource_name,
			a.id AS action_id, a.action_name
		FROM permissions p
		LEFT JOIN resources r ON r.id = p.resource_id
		LEFT JOIN actions a ON a.id = p.action_id
	`

	// Execute the query and scan the results into rows
	err := s.db.Raw(query).Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	// Map for grouping by resource
	resourceMap := make(map[uint]*rolepermdto.GroupedResourcesReponse)

	// Grouping data by resource and actions
	for _, row := range rows {
		if _, exists := resourceMap[row.ResourceID]; !exists {
			// Create new resource entry if it doesn't exist
			resourceMap[row.ResourceID] = &rolepermdto.GroupedResourcesReponse{
				ResourceID:   row.ResourceID,
				ResourceName: row.ResourceName,
				Actions:      []rolepermdto.ActionWithPerms{},
			}
		}

		// Add action with permissions to the resource
		resourceMap[row.ResourceID].Actions = append(resourceMap[row.ResourceID].Actions, rolepermdto.ActionWithPerms{
			ActionID:   row.ActionID,
			ActionName: row.ActionName,
			PermID:     row.PermID,
			PermName:   row.PermName,
			PermCode:   row.PermCode,
		})
	}

	// Convert the map to a slice
	var result []rolepermdto.GroupedResourcesReponse
	for _, resource := range resourceMap {
		result = append(result, *resource)
	}

	return result, nil

}

// Hàm lấy phân quyền theo ID người dùng
func (s *rolePermService) GetPermsByUserID(userID uint) ([]rolepermdto.GetPermByUserIdResult, error) {
	var result []rolepermdto.GetPermByUserIdResult

	// Gọi hàm get_perms_by_user_id
	err := s.db.Raw("SELECT * FROM get_perms_by_user_id(?)", userID).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
