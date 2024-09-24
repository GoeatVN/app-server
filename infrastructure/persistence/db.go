package database

import (
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"food-app/infrastructure/persistence"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Repositories struct {
	User repository.UserRepository
	Food repository.FoodRepository
	db   *gorm.DB
}

func NewDatabase(Driver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	BURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := gorm.Open(Driver, BURL)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User: persistence.NewUserRepository(db),
		Food: persistence.NewFoodRepository(db),
		db:   db,
	}, nil
}

// closes the  database connection
func (s *Repositories) Close() error {
	return s.db.Close()
}

// This migrate all tables
func (s *Repositories) DbMigrate() error {
	return s.db.AutoMigrate(&entity.User{}, &entity.Food{}).Error
}
