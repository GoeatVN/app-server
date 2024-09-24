package persistence

import (
	"fmt"
	"food-app/domain/entity"
	"food-app/domain/repository"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

type Repositories struct {
	User repository.UserRepository
	Food repository.FoodRepository
	db   *gorm.DB
}

func NewDatabase() (*Repositories, error) {
	driver := os.Getenv("DB_DRIVER")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	password := os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")

	dns := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		host,
		port,
		user,
		dbname,
		password)
	db, err := gorm.Open(driver, dns)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &Repositories{
		User: NewUserRepository(db),
		Food: NewFoodRepository(db),
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
