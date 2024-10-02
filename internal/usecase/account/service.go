package account

import (
	"app-server/internal/domain/entity"
	"app-server/internal/persistence/repository"
	"app-server/internal/persistence/repository/postgres"
	"app-server/internal/shared/login"
	"app-server/internal/usecase/auth"
	"fmt"
)

type accountService struct {
	userRepo     *postgres.UserRepository
	userRoleRepo *repository.GenericBaseRepository[entity.UserRole]
	authService  auth.AuthServiceInterface
}

func NewAccountService(repo *postgres.UserRepository, userRoleRepo *repository.GenericBaseRepository[entity.UserRole], authService auth.AuthServiceInterface) ServiceInterface {
	return &accountService{userRepo: repo, userRoleRepo: userRoleRepo, authService: authService}
}

type ServiceInterface interface {
	Login(loginDto login.LoginRequest) (*login.LoginResponse, error)
}

func (s *accountService) Login(loginDto login.LoginRequest) (*login.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(loginDto.Username)
	if err != nil {
		return nil, err
	}

	// Compare hashed password with plain text password
	if s.authService.CheckPassword(user.Password, loginDto.Password) != nil {
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
	token, err := s.authService.GenerateJWT(user.ID, roleIDs, user.Username)
	if err != nil {
		return nil, err
	}

	return &login.LoginResponse{AccessToken: token}, nil
}
