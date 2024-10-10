package repository

import (
	"fmt"
	"gorm.io/gorm"
)

// GenericBaseRepository là một repository tổng quát cho mọi entity/model
type GenericBaseRepository[T any] struct {
	db *gorm.DB
}

// NewGenericBaseRepository tạo mới một repository với GORM DB
func NewGenericBaseRepository[T any](db *gorm.DB) *GenericBaseRepository[T] {
	return &GenericBaseRepository[T]{db: db}
}

// Create thêm một bản ghi mới vào database
func (r *GenericBaseRepository[T]) Create(entity *T) error {
	var currentSearchPath string
	r.db.Raw("SHOW search_path").Scan(&currentSearchPath)
	fmt.Println("Current search path:", currentSearchPath)
	return r.db.Create(entity).Error
}

// Create thêm nhiều bản ghi mới vào database
func (r *GenericBaseRepository[T]) CreateMany(entities []T) error {
	return r.db.Create(entities).Error
}

// Update cập nhật một bản ghi đã tồn tại
func (r *GenericBaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Cập nhật tất cả các trường của bản ghi
func (r *GenericBaseRepository[T]) UpdateAll(entity *T) error {
	return r.db.Model(entity).Updates(entity).Error
}

// Delete xóa một bản ghi theo ID
func (r *GenericBaseRepository[T]) Delete(id uint) error {
	return r.db.Delete(new(T), id).Error
}

// FindAll trả về tất cả các bản ghi của model
func (r *GenericBaseRepository[T]) FindAll() ([]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	return entities, err
}

// FindByID trả về một bản ghi theo ID
func (r *GenericBaseRepository[T]) FindByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// Where trả về một DB với điều kiện query
func (r *GenericBaseRepository[T]) Where(query interface{}, args ...interface{}) *gorm.DB {
	return r.db.Where(query, args...)
}

// First trả về bản ghi đầu tiên tìm thấy
func (r *GenericBaseRepository[T]) First(entity *T) error {
	return r.db.First(entity).Error
}

// Count trả về số lượng bản ghi
func (r *GenericBaseRepository[T]) Count() int64 {
	var entity T
	var count int64
	r.db.Model(entity).Count(&count)
	return count
}

// Kiểm tra tồn tại theo điều kiện query
func (r *GenericBaseRepository[T]) Exists(query interface{}, args ...interface{}) bool {
	var entity T
	err := r.db.Where(query, args...).First(entity).Error
	return err == nil || err != gorm.ErrRecordNotFound

}
