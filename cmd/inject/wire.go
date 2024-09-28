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
	"app-server/internal/usecase/user"
	"github.com/google/wire"
	"gorm.io/gorm"
)

// create provider function for GenericBaseRepository
// provideGenericBaseRepository sẽ inject database connection và trả về GenericBaseRepository
func provideGenericBaseRepository(db *gorm.DB) *repository.GenericBaseRepository[entity.UserRole] {
	return repository.NewGenericBaseRepository[entity.UserRole](db)
}

// InitializeServer sẽ inject tất cả các dependencies và trả về Server
func InitializeServer(config *config.Config) (*server.HTTPServer, error) {
	wire.Build(
		database.Connect,           // Inject database connection
		postgres.NewUserRepository, // Inject UserRepository
		provideGenericBaseRepository,
		user.NewService,           // Inject UserService
		auth.NewAuthService,       // Inject AuthService
		account.NewAccountService, // Inject AccountService
		v1.NewUserHandler,         // Inject UserHandler
		v1.NewAccountHandler,
		server.NewHTTPServer,
	)
	return &server.HTTPServer{}, nil
}
