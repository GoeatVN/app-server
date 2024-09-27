package main

import (
	"app-server/internal/domain/entity"
	"app-server/internal/infrastructure/config"
	"app-server/internal/infrastructure/database"
	"app-server/internal/infrastructure/server"
	"app-server/internal/interface/api/handler/v1"
	"app-server/internal/persistence/repository"
	"app-server/internal/persistence/repository/postgres"
	"app-server/internal/usecase/account"
	"app-server/internal/usecase/user"
	"log"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	//server, err := InitializeServer(config)
	db, err := database.Connect(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepository := postgres.NewUserRepository(db)
	userServiceInterface := user.NewService(userRepository)
	userHandler := v1.NewUserHandler(userServiceInterface)

	// Account-related dependencies
	accountRepository := repository.NewGenericBaseRepository[entity.UserRole](db)
	accountServiceInterface := account.NewAccountService(userRepository, accountRepository)
	accountHandler := v1.NewAccountHandler(accountServiceInterface)

	server := server.NewHTTPServer(config, userHandler, accountHandler)

	if err != nil {
		log.Fatalf("Failed to initialize API: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
