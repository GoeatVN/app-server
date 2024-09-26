//go:build wireinject
// +build wireinject

package main

import (
	"app-server/internal/infrastructure/config"
	"app-server/internal/infrastructure/database"
	"app-server/internal/interface/api/handler/v1"

	"app-server/internal/persistence/repository/postgres"
	"app-server/internal/usecase/user"

	"app-server/internal/infrastructure/server"

	"github.com/google/wire"
)

// InitializeServer sẽ inject tất cả các dependencies và trả về Server
func InitializeServer(config *config.Config) (*server.HTTPServer, error) {
	wire.Build(
		database.Connect,           // Inject database connection
		postgres.NewUserRepository, // Inject UserRepository
		user.NewService,            // Inject UserService
		v1.NewUserHandler,          // Inject UserHandler
		server.NewHTTPServer,
	)
	return &server.HTTPServer{}, nil
}
