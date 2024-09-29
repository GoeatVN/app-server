package user

import (
	"app-server/internal/domain/entity"
	"app-server/internal/persistence/repository/postgres"
	"app-server/internal/shared/userdto"
	"app-server/internal/usecase/auth"
	"fmt"
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
}

func NewService(repo *postgres.UserRepository, authService auth.AuthServiceInterface) ServiceInterface {
	return &service{repo: repo, authService: authService}
}

func (s *service) GetAllUsers() ([]entity.User, error) {
	return s.repo.FindAll()
}

func (s *service) GetUserByID(id uint) (*entity.User, error) {
	fmt.Println("service.GetUserByID")
	return s.repo.FindByID(id)
}

func (s *service) CreateUser(request *userdto.AddUserRequest) (string, error) {
	// Mã hóa mật khẩu
	hashedPassword, err := s.authService.HashPassword(request.Password)
	if err != nil {
		return "", err // Trả về lỗi nếu mã hóa thất bại
	}
	// Tạo đối tượng User mới
	newUser := &entity.User{
		Username: request.Username,
		Password: hashedPassword,
		Email:    request.Email,
		Phone:    request.Phone,
	}

	errCreate := s.repo.Create(newUser)
	if errCreate != nil {
		return "", errCreate
	}

	return "Tạo người dùng thành công", nil
}

func (s *service) UpdateUser(user *entity.User) error {
	return s.repo.Update(user)
}

func (s *service) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
