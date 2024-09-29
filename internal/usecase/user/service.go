package user

import (
	"app-server/internal/domain/entity"
	"app-server/internal/persistence/repository/postgres"
	"app-server/internal/shared/userdto"
	"app-server/internal/usecase/auth"
	"fmt"
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
}

func NewService(repo *postgres.UserRepository, authService auth.AuthServiceInterface, db *gorm.DB) ServiceInterface {
	return &service{repo: repo, authService: authService, db: db}
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

	type Result struct {
		IsSuccess    bool
		ErrorCode    int
		ErrorMessage string
	}

	var isSuccess bool
	var errorCode int
	var errorMessage string

	// Execute the CALL statement
	err = s.db.Exec(`
        CALL demovcs.create_user(?, ?, ?, ?, ?, ?, ?, ?, ?);
    `,
		request.Username,
		hashedPassword,
		request.Email,
		request.Phone,
		request.Fullname,
		request.CreatedBy,
		gorm.Expr("?", &isSuccess),
		gorm.Expr("?", &errorCode),
		gorm.Expr("?", &errorMessage),
	).Error

	// Tạo đối tượng kết quả
	type CreateUserResult struct {
		IsSuccess    bool
		ErrorCode    int
		ErrorMessage string
	}
	result := CreateUserResult{
		IsSuccess:    isSuccess,
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}

	// In kết quả
	fmt.Printf("Is Success: %v\n", result.IsSuccess)
	fmt.Printf("Error Code: %d\n", result.ErrorCode)
	fmt.Printf("Error Message: %s\n", result.ErrorMessage)

	// Lấy giá trị của các tham số OUT
	return "User created successfully", nil
}

func (s *service) UpdateUser(user *entity.User) error {
	return s.repo.Update(user)
}

func (s *service) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
