package auth

import (
	"app-server/internal/domain/entity"
	"app-server/internal/infrastructure/config"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type authService struct {
	config config.Config
}

func NewAuthService(config *config.Config) AuthServiceInterface {
	return &authService{config: *config}
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	RoleIDs  []uint `json:"role_ids"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthServiceInterface interface {
	VerifyToken(tokenString string) (*entity.AuthClaims, error)
	GenerateJWT(userID uint, roleIDs []uint, username string) (string, error)
	GetClaims(tokenString string) (*entity.AuthClaims, error)
	HashPassword(password string) (string, error)
	CheckPassword(hashedPassword, password string) error
}

// VerifyToken verifies the JWT token
func (s *authService) VerifyToken(tokenString string) (*entity.AuthClaims, error) {
	// Parse the JWT token with custom claims
	token, err := jwt.ParseWithClaims(tokenString, &entity.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	// Kiểm tra token có hợp lệ không và có chứa claims không
	if claims, ok := token.Claims.(*entity.AuthClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token or claims")
	}
}

func (s *authService) GenerateJWT(userID uint, roleIDs []uint, username string) (string, error) {

	// Set token expiration time
	expirationTime := time.Now().Add(time.Minute * time.Duration(s.config.JWT.TokenExpiry))
	// Create the JWT claims, which includes the user ID, role IDs, and username
	claims := &entity.AuthClaims{
		UserID:   userID,
		RoleIDs:  roleIDs,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) GetClaims(tokenString string) (*entity.AuthClaims, error) {
	// Parse the JWT token
	token, err := jwt.ParseWithClaims(tokenString, &entity.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	// Get the claims
	claims, ok := token.Claims.(*entity.AuthClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}

// Hàm HashPassword để mã hóa mật khẩu
func (s *authService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword kiểm tra xem mật khẩu nhập vào có khớp với mật khẩu đã mã hóa hay không
func (s *authService) CheckPassword(hashedPassword, password string) error {
	// So sánh mật khẩu đã nhập với mật khẩu đã mã hóa
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Nếu không khớp, trả về lỗi
		return errors.New("mật khẩu không chính xác")
	}
	// Nếu khớp, trả về nil (không có lỗi)
	return nil
}
