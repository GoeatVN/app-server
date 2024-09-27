package account

import (
	"app-server/internal/domain/entity"
	"app-server/internal/persistence/repository"
	"app-server/internal/persistence/repository/postgres"
	"app-server/internal/shared/login"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type accountService struct {
	userRepo     *postgres.UserRepository
	userRoleRepo *repository.GenericBaseRepository[entity.UserRole]
}

func NewAccountService(repo *postgres.UserRepository, userRoleRepo *repository.GenericBaseRepository[entity.UserRole]) ServiceInterface {
	return &accountService{userRepo: repo, userRoleRepo: userRoleRepo}
}

type ServiceInterface interface {
	Login(loginDto login.LoginDTO) (*string, error)
}

func (s *accountService) Login(loginDto login.LoginDTO) (*string, error) {
	user, err := s.userRepo.FindByUsername(loginDto.Username)
	if err != nil {
		return nil, err
	}

	// Compare hashed password with plain text password
	if user.Password != loginDto.Password {
		return nil, fmt.Errorf("invalid credentials")
	}

	var userRole []entity.UserRole
	if err := s.userRoleRepo.Where("user_id = ?", user.ID).Find(&userRole).Error; err != nil {
		return nil, err
	}
	var roleIDs []uint
	for _, item := range userRole {
		roleIDs = append(roleIDs, item.RoleID)
	}

	// Generate JWT
	secretKey := "your_secret_key" // Replace with your actual secret key
	token, err := GenerateJWT(user.ID, roleIDs, user.Username, secretKey)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func GenerateJWT(userID uint, roleIDs []uint, username string, secretKey string) (string, error) {
	type Claims struct {
		UserID   uint   `json:"user_id"`
		RoleIDs  []uint `json:"role_ids"`
		Username string `json:"username"`
		jwt.StandardClaims
	}
	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes the user ID, role IDs, and username
	claims := &Claims{
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
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
