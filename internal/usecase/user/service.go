package user

import (
	"app-server/internal/domain/entity"
	"app-server/internal/persistence/repository/postgres"
	"app-server/internal/shared/userdto"
	"app-server/internal/usecase/auth"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/gorm"
)

type ServiceInterface interface {
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	CreateUser(request *userdto.AddUserRequest) (string, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uint) error
}

type service struct {
	repo        *postgres.UserRepository
	authService auth.AuthServiceInterface
	db          *gorm.DB
	dbPool      *pgxpool.Pool
}

func NewService(repo *postgres.UserRepository, authService auth.AuthServiceInterface, db *gorm.DB, dbPool *pgxpool.Pool) ServiceInterface {
	return &service{repo: repo, authService: authService, db: db, dbPool: dbPool}
}

func (s *service) GetAllUsers() ([]entity.User, error) {
	return s.repo.FindAll()
}

func (s *service) GetUserByID(id uint) (*entity.User, error) {
	fmt.Println("service.GetUserByID")
	return s.repo.FindByID(id)
}

func (s *service) CreateUser(request *userdto.AddUserRequest) (string, error) {
	// Hash the password
	hashedPassword, err := s.authService.HashPassword(request.Password)
	if err != nil {
		return "", err // Return error if hashing fails
	}

	// Define the OUT parameters
	var isSuccess bool
	var errorCode int
	var errorMessage string

	// Call the stored procedure using QueryRow
	err = s.dbPool.QueryRow(
		context.Background(),
		"CALL create_user($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		request.Username,
		hashedPassword,
		request.Email,
		request.Phone,
		request.FullName,
		request.CreatedBy,
		&isSuccess,
		&errorCode,
		&errorMessage,
	).Scan(&isSuccess, &errorCode, &errorMessage)

	if err != nil {
		return "", fmt.Errorf("failed to execute procedure: %v", err)
	}

	// Check the result
	if !isSuccess {
		return "", fmt.Errorf("Error Code: %d, Message: %s", errorCode, errorMessage)
	}

	return "User created successfully", nil
}
func (s *service) UpdateUser(user *entity.User) error {
	return s.repo.Update(user)
}

func (s *service) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
