package server

import (
	"app-server/internal/domain/enum"
	"app-server/internal/infrastructure/config"
	"app-server/internal/infrastructure/middleware"
	"app-server/internal/interface/api/handler/v1"
	"app-server/internal/usecase/auth"
	"app-server/internal/usecase/rolepermission"
	"github.com/gin-gonic/gin"
)

// Update the HTTPServer struct
type HTTPServer struct {
	router          *gin.Engine
	config          *config.Config
	userHandler     *v1.UserHandler
	accountHandler  *v1.AccountHandler
	authService     auth.AuthServiceInterface
	rolePermService rolepermission.RolePermServiceInterface
	rolePermHandler *v1.RolePermHandler
	systemHandler   *v1.SystemHandler // Add SystemHandler
}

// Update the NewHTTPServer function
func NewHTTPServer(
	config *config.Config,
	userHandler *v1.UserHandler,
	accountHandler *v1.AccountHandler,
	rolePermHandler *v1.RolePermHandler,
	authService auth.AuthServiceInterface,
	rolePermService rolepermission.RolePermServiceInterface,
	systemHandler *v1.SystemHandler, // Add SystemHandler
) *HTTPServer {
	router := gin.Default()

	// Register middleware
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RateLimiterMiddleware())
	router.Use(middleware.CORS())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.ResponseHandlerMiddleware())

	server := &HTTPServer{
		router:          router,
		config:          config,
		userHandler:     userHandler,
		accountHandler:  accountHandler,
		authService:     authService,
		rolePermService: rolePermService,
		rolePermHandler: rolePermHandler,
		systemHandler:   systemHandler, // Initialize SystemHandler
	}

	server.setupRoutes()

	return server
}

// Update the setupRoutes method
func (s *HTTPServer) setupRoutes() {
	// Public routes
	s.router.POST("/api/account/login", s.accountHandler.Login)

	// Create middleware auth
	authMiddleware := middleware.NewAuthMiddleware(s.authService, s.rolePermService)

	api := s.router.Group("/api")
	{
		api.Use(authMiddleware.AuthN())

		api.GET("/users", authMiddleware.AuthZ(enum.Resource.User, enum.Action.View), s.userHandler.GetAllUsers)
		api.GET("/users/:id", authMiddleware.AuthZ(enum.Resource.User, enum.Action.View), s.userHandler.GetUserByID)
		api.POST("/users/add", authMiddleware.AuthZ(enum.Resource.User, enum.Action.Add), s.userHandler.CreateUser)
		api.POST("/users/:id/modify", authMiddleware.AuthZ(enum.Resource.User, enum.Action.Update), s.userHandler.UpdateUser)
		api.GET("/users/:id/perms", authMiddleware.AuthZ(enum.Resource.User, enum.Action.View), s.rolePermHandler.GetPermsByUserID)

		api.GET("/resources", authMiddleware.AuthZ(enum.Resource.Role, enum.Action.View), s.rolePermHandler.GetResources)
		api.GET("/roles", authMiddleware.AuthZ(enum.Resource.Role, enum.Action.View), s.rolePermHandler.GetAllRolePerms)
		api.GET("/roles/:id", authMiddleware.AuthZ(enum.Resource.Role, enum.Action.View), s.rolePermHandler.GetRolePermsById)
		api.POST("/roles/add", authMiddleware.AuthZ(enum.Resource.Role, enum.Action.Add, enum.Action.Update), s.rolePermHandler.AddNewRole)
		api.POST("/roles/:id/modify", authMiddleware.AuthZ(enum.Resource.Role, enum.Action.Update), s.rolePermHandler.ModifyRole)
		api.POST("/roles/asign-role", authMiddleware.AuthZ(enum.Resource.Role, enum.Action.Add, enum.Action.Update), s.rolePermHandler.AssignRoleToUser)

		// Add the new route for LoadComboboxData
		api.POST("/combobox/load", s.systemHandler.LoadComboboxDataHandler)
	}
}
