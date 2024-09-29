// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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
	"gorm.io/gorm"
)

// Injectors from wire.go:

// InitializeServer injects all dependencies and returns Server
func InitializeServer(config2 *config.Config) (*server.HTTPServer, error) {
	db, err := database.Connect(config2)
	if err != nil {
		return nil, err
	}
	userRepository := postgres.NewUserRepository(db)
	authServiceInterface := auth.NewAuthService(config2)
	serviceInterface := user.NewService(userRepository, authServiceInterface, db)
	userHandler := v1.NewUserHandler(serviceInterface)
	genericBaseRepository := provideUserRoleRepo(db)
	accountServiceInterface := account.NewAccountService(userRepository, genericBaseRepository, authServiceInterface)
	accountHandler := v1.NewAccountHandler(accountServiceInterface)
	repositoryGenericBaseRepository := provideRoleRepo(db)
	genericBaseRepository2 := provideRolePermissionRepo(db)
	genericBaseRepository3 := providePermissionRepo(db)
	genericBaseRepository4 := provideResourceRepo(db)
	genericBaseRepository5 := provideActionRepo(db)
	rolePermServiceInterface := rolepermission.NewRolePermService(userRepository, genericBaseRepository, repositoryGenericBaseRepository, genericBaseRepository2, genericBaseRepository3, genericBaseRepository4, genericBaseRepository5, db)
	rolePermHandler := v1.NewRolePermHandler(rolePermServiceInterface)
	httpServer := server.NewHTTPServer(config2, userHandler, accountHandler, rolePermHandler, authServiceInterface, rolePermServiceInterface)
	return httpServer, nil
}

// wire.go:

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
