package middleware

import (
	"app-server/internal/domain/enum"
	"app-server/internal/usecase/auth"
	"app-server/internal/usecase/rolepermission"
	"app-server/internal/utils"
	"app-server/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	AuthService     auth.AuthServiceInterface
	RolePermService rolepermission.RolePermServiceInterface
}

func NewAuthMiddleware(authService auth.AuthServiceInterface, rolePermService rolepermission.RolePermServiceInterface) *AuthMiddleware {
	return &AuthMiddleware{
		AuthService:     authService,
		RolePermService: rolePermService,
	}
}

// Authentication kiểm tra Bearer token từ header HTTP
func (s *AuthMiddleware) AuthN() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		// Kiểm tra header Authorization có chứa Bearer token không
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing or invalid token")
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Kiểm tra token bằng cách sử dụng authService
		claims, err := s.AuthService.VerifyToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token")
			c.Abort()
			return
		}
		// Lưu claims vào Gin context để dùng sau
		c.Set("tokenClaims", claims)
		c.Next()
	}
}

// Authorization kiểm tra quyền truy cập dựa trên vai trò người dùng
func (s *AuthMiddleware) AuthZ(resource enum.ResourceCode, actions ...enum.ActionCode) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy claims từ Gin context (đã lưu trong AuthN)
		claim, errGetClaim := utils.GetAuthClaims(c)
		if errGetClaim != nil {
			response.Error(c, errGetClaim.HTTPCode, errGetClaim.Code, errGetClaim.Message)
			c.Abort()
			return
		}

		perms, err := s.RolePermService.GetPermsByUserID(claim.UserID)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, "ERROR", "Could not get permissions")
			c.Abort()
			return
		}

		hasPermission := false
		for _, perm := range perms {
			for _, action := range actions {
				if perm.ResourceCode == string(resource) && perm.ActionCode == string(action) {
					hasPermission = true
					break
				}
			}
			if hasPermission {
				break
			}
		}

		if !hasPermission {
			response.Error(c, http.StatusForbidden, "FORBIDDEN", "You do not have permission to access this resource")
			c.Abort()
			return
		}

		// Nếu vai trò hợp lệ, tiếp tục xử lý request
		c.Next()
	}
}
