package server

import (
	"app-server/internal/infrastructure/config"
	"app-server/internal/infrastructure/middleware"
	"app-server/internal/interface/api/handler/v1"
	"app-server/internal/usecase/auth"
	"fmt"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct {
	router         *gin.Engine
	config         *config.Config
	userHandler    *v1.UserHandler
	accountHandler *v1.AccountHandler
	authService    auth.AuthServiceInterface
}

func NewHTTPServer(
	config *config.Config,
	userHandler *v1.UserHandler,
	accountHandler *v1.AccountHandler,
	authService auth.AuthServiceInterface,
) *HTTPServer {
	router := gin.Default()

	// Đăng ký các middleware
	router.Use(middleware.LoggerMiddleware()) // Ghi log
	//router.Use(middleware.AuthMiddleware())            // Xác thực token
	router.Use(middleware.RateLimiterMiddleware())     // Giới hạn số lượng yêu cầu từ một IP
	router.Use(middleware.ErrorHandler())              // Xử lý lỗi phát sinh
	router.Use(middleware.ResponseHandlerMiddleware()) // Chuẩn hóa kết quả trả về

	// Áp dụng middleware Authorization để giới hạn quyền truy cập cho vai trò "admin"
	// adminGroup := router.Group("/admin")
	// adminGroup.Use(middleware.AuthorizationMiddleware("admin"))
	// adminGroup.GET("/users", s.userHandler.GetUsers)

	// // Route không cần kiểm tra quyền, mọi người dùng đều truy cập được
	// router.POST("/users", s.userHandler.CreateUser)

	server := &HTTPServer{
		router:         router,
		config:         config,
		userHandler:    userHandler,
		accountHandler: accountHandler,
		authService:    authService,
	}

	server.setupRoutes()

	return server
}

func (s *HTTPServer) setupRoutes() {
	// Route không cần kiểm tra quyền, mọi người dùng đều truy cập được
	s.router.POST("/api/account/login", s.accountHandler.Login)

	api := s.router.Group("/api")
	{
		api.Use(middleware.AuthenticationMiddleware(s.authService))

		api.POST("/users", s.userHandler.CreateUser)
		api.GET("/users/", s.userHandler.GetUsers)
		api.GET("/users/:id", s.userHandler.GetUserByID)
		// Add other routes as needed
	}
}

func (s *HTTPServer) Run() error {
	return s.router.Run(fmt.Sprintf(":%d", s.config.App.Port))
}
