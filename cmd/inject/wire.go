//go:build wireinject
// +build wireinject

package inject

import (
	"app-server/internal/domain/entity"
	"app-server/internal/infrastructure/config"
	"app-server/internal/infrastructure/database"
	"app-server/internal/infrastructure/server"
	"app-server/internal/interface/api/handler/v1"
	"app-server/internal/persistence/repository"
	"app-server/internal/persistence/repository/postgres"
	"app-server/internal/usecase/account"
	"app-server/internal/usecase/auth"
	"app-server/internal/usecase/rolepermission"
	"app-server/internal/usecase/user"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// provideUserRoleRepo injects the database connection and returns UserRoleRepo
func provideUserRoleRepo(db *gorm.DB) *repository.GenericBaseRepository[entity.UserRole] {
	return repository.NewGenericBaseRepository[entity.UserRole](db)
}

// provideRoleRepo injects the database connection and returns RoleRepo
func provideRoleRepo(db *gorm.DB) *repository.GenericBaseRepository[entity.Role] {
	return repository.NewGenericBaseRepository[entity.Role](db)
}

// provideRolePermissionRepo injects the database connection and returns RolePermissionRepo
func provideRolePermissionRepo(db *gorm.DB) *repository.GenericBaseRepository[entity.RolePermission] {
	return repository.NewGenericBaseRepository[entity.RolePermission](db)
}

// providePermissionRepo injects the database connection and returns PermissionRepo
func providePermissionRepo(db *gorm.DB) *repository.GenericBaseRepository[entity.Permission] {
	return repository.NewGenericBaseRepository[entity.Permission](db)
}

// provideResourceRepo injects the database connection and returns ResourceRepo
func provideResourceRepo(db *gorm.DB) *repository.GenericBaseRepository[entity.Resource] {
	return repository.NewGenericBaseRepository[entity.Resource](db)
}

// provideActionRepo injects the database connection and returns ActionRepo
func provideActionRepo(db *gorm.DB) *repository.GenericBaseRepository[entity.Action] {
	return repository.NewGenericBaseRepository[entity.Action](db)
}

// InitializeServer injects all dependencies and returns Server
func InitializeServer(config *config.Config) (*server.HTTPServer, error) {
	wire.Build(
		database.Connect, // Inject database connection
		database.ConnectToDBPool,
		postgres.NewUserRepository, // Inject UserRepository
		provideUserRoleRepo,
		provideRoleRepo,
		provideRolePermissionRepo,
		providePermissionRepo,
		provideResourceRepo,
		provideActionRepo,
		user.NewService,                   // Inject UserService
		auth.NewAuthService,               // Inject AuthService
		account.NewAccountService,         // Inject AccountService
		rolepermission.NewRolePermService, // Inject RolePermService
		soil_analysis.NewSoilAnalysisServiceInterface,
		v1.NewUserHandler, // Inject UserHandler
		system.NewComboboxService,         // Inject SystemService
		v1.NewUserHandler,                 // Inject UserHandler
		v1.NewAccountHandler,
		v1.NewRolePermHandler,
		v1.NewSystemHandler, // Inject SystemHandler
		v1.NewSoilAnalysisHandler,
		server.NewHTTPServer,
	)
	return &server.HTTPServer{}, nil
}
