package server

import (
	"app-server/internal/domain/enum"
	"app-server/internal/infrastructure/config"
	"app-server/internal/infrastructure/middleware"
	v1 "app-server/internal/interface/api/handler/v1"
	"app-server/internal/usecase/auth"
	"app-server/internal/usecase/rolepermission"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Update the HTTPServer struct
type HTTPServer struct {
	router              *gin.Engine
	config              *config.Config
	userHandler         *v1.UserHandler
	accountHandler      *v1.AccountHandler
	rolePermHandler     *v1.RolePermHandler
	soilAnalysisHandler *v1.SoilAnalysisHandler
	authService         auth.AuthServiceInterface
	rolePermService     rolepermission.RolePermServiceInterface
	systemHandler   *v1.SystemHandler // Add SystemHandler
}

type Handlers struct {
	UserHandler         *v1.UserHandler
	AccountHandler      *v1.AccountHandler
	RolePermHandler     *v1.RolePermHandler
	SoilAnalysisHandler *v1.SoilAnalysisHandler
}

type Services struct {
	AuthService     auth.AuthServiceInterface
	RolePermService rolepermission.RolePermServiceInterface
}

// Update the NewHTTPServer function
func NewHTTPServer(
	config *config.Config,
	UserHandler *v1.UserHandler,
	AccountHandler *v1.AccountHandler,
	RolePermHandler *v1.RolePermHandler,
	SoilAnalysisHandler *v1.SoilAnalysisHandler,
	AuthService auth.AuthServiceInterface,
	RolePermService rolepermission.RolePermServiceInterface,
	// redisCache *cache.RedisCache,
	SystemHandler *v1.SystemHandler, // Add SystemHandler
) *HTTPServer {
	router := gin.Default()

	// Đăng ký các middleware
	router.Use(middleware.LoggerMiddleware()) // Ghi log
	//router.Use(middleware.AuthMiddleware())            // Xác thực token
	router.Use(middleware.RateLimiterMiddleware())     // Giới hạn số lượng yêu cầu từ một IP
	router.Use(middleware.CORS())                      // Xử lý CORS
	router.Use(middleware.ErrorHandler())              // Xử lý lỗi phát sinh
	router.Use(middleware.ResponseHandlerMiddleware()) // Chuẩn hóa kết quả trả về
	// Đăng ký middleware caching và chuẩn hóa kết quả trả về
	//router.Use(middleware.CachingMiddleware(s.redisCache, 10*time.Minute)) // Cache trong 10 phút

	server := &HTTPServer{
		router:              router,
		config:              config,
		userHandler:         UserHandler,
		accountHandler:      AccountHandler,
		authService:         AuthService,
		rolePermService:     RolePermService,
		rolePermHandler:     RolePermHandler,
		soilAnalysisHandler: SoilAnalysisHandler,
		systemHandler:   SystemHandler,
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
	// Route không cần kiểm tra quyền, mọi người dùng đều truy cập được
	apiAnonymos := s.router.Group("/api")
	{
		apiAnonymos.POST("/soil-analysis", s.soilAnalysisHandler.SaveSoilAnalysis)
	}
}

func (s *HTTPServer) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.App.Port))
}
