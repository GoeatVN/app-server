package application

import (
	"food-app/domain/entity"
	"food-app/domain/repository"
)

// userApp struct holds methods for user use cases
// contructor
type userApp struct {
	us repository.UserRepository
}

// UserApp implements the UserAppInterface
var _ UserAppInterface = &userApp{}

// UserAppInterface defines the methods that any implementation of user application logic must provide.
// It includes methods for saving a user, retrieving all users, retrieving a user by ID, and retrieving a user by email and password.
type UserAppInterface interface {
	SaveUser(*entity.User) (*entity.User, map[string]string)
	GetUsers() ([]entity.User, error)
	GetUser(uint64) (*entity.User, error)
	GetUserByEmailAndPassword(*entity.User) (*entity.User, map[string]string)
}

func (u *userApp) SaveUser(user *entity.User) (*entity.User, map[string]string) {
	return u.us.SaveUser(user)
}

func (u *userApp) GetUser(userId uint64) (*entity.User, error) {
	return u.us.GetUser(userId)
}

func (u *userApp) GetUsers() ([]entity.User, error) {

	return u.us.GetUsers()
}

func (u *userApp) GetUserByEmailAndPassword(user *entity.User) (*entity.User, map[string]string) {
	return u.us.GetUserByEmailAndPassword(user)
}
