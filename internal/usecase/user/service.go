package user

import (
	"app-server/internal/domain/entity"
	"app-server/internal/persistence/repository/postgres"
	"fmt"
)

type ServiceInterface interface {
	GetAllUsers() ([]entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	DeleteUser(id uint) error
}

type service struct {
	repo *postgres.UserRepository
}

func NewService(repo *postgres.UserRepository) ServiceInterface {
	return &service{repo: repo}
}

func (s *service) GetAllUsers() ([]entity.User, error) {
	return s.repo.FindAll()
}

func (s *service) GetUserByID(id uint) (*entity.User, error) {
	fmt.Println("service.GetUserByID")
	return s.repo.FindByID(id)
}

func (s *service) CreateUser(user *entity.User) error {
	return s.repo.Create(user)
}

func (s *service) UpdateUser(user *entity.User) error {
	return s.repo.Update(user)
}

func (s *service) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}
