package postgres

import (
	"app-server/internal/domain/entity"
	"app-server/internal/persistence/repository"

	"gorm.io/gorm"
)

// UserRepository sử dụng GenericBaseRepository thông qua embedding
type UserRepository struct {
	*repository.GenericBaseRepository[entity.User]
}

// NewUserRepository khởi tạo repository cho User
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		GenericBaseRepository: repository.NewGenericBaseRepository[entity.User](db),
	}
}

// Hàm riêng cho User nếu cần
func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
