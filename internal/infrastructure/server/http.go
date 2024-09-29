package server

import (
	"app-server/internal/domain/enum"
	"app-server/internal/infrastructure/config"
	"app-server/internal/infrastructure/middleware"
	"app-server/internal/interface/api/handler/v1"
	"app-server/internal/usecase/auth"
	"app-server/internal/usecase/rolepermission"
	"fmt"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	router          *gin.Engine
	config          *config.Config
	userHandler     *v1.UserHandler
	accountHandler  *v1.AccountHandler
	authService     auth.AuthServiceInterface
	rolePermService rolepermission.RolePermServiceInterface
	rolePermHandler *v1.RolePermHandler
}

func NewHTTPServer(
	config *config.Config,
	userHandler *v1.UserHandler,
	accountHandler *v1.AccountHandler,
	rolePermHandler *v1.RolePermHandler,
	authService auth.AuthServiceInterface,
	rolePermService rolepermission.RolePermServiceInterface,
	// redisCache *cache.RedisCache,
) *HTTPServer {
	router := gin.Default()

	// Đăng ký các middleware
	router.Use(middleware.LoggerMiddleware()) // Ghi log
	//router.Use(middleware.AuthMiddleware())            // Xác thực token
	router.Use(middleware.RateLimiterMiddleware())     // Giới hạn số lượng yêu cầu từ một IP
	router.Use(middleware.ErrorHandler())              // Xử lý lỗi phát sinh
	router.Use(middleware.ResponseHandlerMiddleware()) // Chuẩn hóa kết quả trả về
	// Đăng ký middleware caching và chuẩn hóa kết quả trả về
	//router.Use(middleware.CachingMiddleware(s.redisCache, 10*time.Minute)) // Cache trong 10 phút

	// Áp dụng middleware Authorization để giới hạn quyền truy cập cho vai trò "admin"
	// adminGroup := router.Group("/admin")
	// adminGroup.Use(middleware.AuthorizationMiddleware("admin"))
	// adminGroup.GET("/users", s.userHandler.GetUsers)

	// // Route không cần kiểm tra quyền, mọi người dùng đều truy cập được
	// router.POST("/users", s.userHandler.CreateUser)

	server := &HTTPServer{
		router:          router,
		config:          config,
		userHandler:     userHandler,
		accountHandler:  accountHandler,
		authService:     authService,
		rolePermService: rolePermService,
		rolePermHandler: rolePermHandler,
	}

	server.setupRoutes()

	return server
}

func (s *HTTPServer) setupRoutes() {
	// Route không cần kiểm tra quyền, mọi người dùng đều truy cập được
	s.router.POST("/api/account/login", s.accountHandler.Login)

	// Create middleware auth
	authMiddleware := middleware.NewAuthMiddleware(s.authService, s.rolePermService)

	api := s.router.Group("/api")
	{
		api.Use(authMiddleware.AuthN())

		api.POST("/users", authMiddleware.AuthZ(enum.Resource.User, enum.Action.Create), s.userHandler.CreateUser)
		api.GET("/users/", s.userHandler.GetUsers)
		api.GET("/users/:id", s.userHandler.GetUserByID)

		api.POST("/roles", s.rolePermHandler.AddNewRole)
		api.POST("/roles/modify", s.rolePermHandler.ModifyRole)
		api.POST("/roles/asign-role", s.rolePermHandler.AssignRoleToUser)
		api.GET("/role-perm", s.rolePermHandler.GetAllRolePerms)
		api.GET("/role-perm/:id", s.rolePermHandler.GetRolePermsById)
		api.GET("/role-perm/group-resource", authMiddleware.AuthZ(enum.Resource.Role, enum.Action.View), s.rolePermHandler.GetGroupResources)
		api.GET("/user-perm/:id", authMiddleware.AuthZ(enum.Resource.Role, enum.Action.View), s.rolePermHandler.GetPermsByUserID)
		// Add other routes as needed
	}
}

func (s *HTTPServer) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.App.Port))
}
