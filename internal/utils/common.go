package utils

import (
	"app-server/internal/domain/entity"
	"app-server/internal/shared/common"
	"app-server/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetAuthClaims retrieves AuthClaims from the Gin context
func GetAuthClaims(c *gin.Context) (*entity.AuthClaims, *common.GetAuthClaimsError) {
	// Lấy claims từ Gin context (đã lưu trong AuthN)
	claimCtx, exists := c.Get("tokenClaims")
	if !exists {
		// Nếu không tồn tại claims, trả về lỗi Unauthorized
		//response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Missing claims")
		//c.Abort()

		return nil, &common.GetAuthClaimsError{http.StatusUnauthorized, "UNAUTHORIZED", "Missing claims"}
	}

	// Chuyển đổi kiểu dữ liệu (nếu cần)
	claim, ok := claimCtx.(*entity.AuthClaims)
	if !ok {
		// Nếu không thể chuyển đổi, trả về lỗi Unauthorized
		response.Error(c, http.StatusUnauthorized, "UNAUTHORIZED", "Invalid claims")
		c.Abort()
		return nil, &common.GetAuthClaimsError{http.StatusUnauthorized, "UNAUTHORIZED", "Missing claims"}
	}

	// Trả về AuthClaims đã lấy từ context
	return claim, nil
}
