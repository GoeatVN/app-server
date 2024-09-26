package database

import (
	"app-server/internal/infrastructure/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect mở kết nối tới PostgreSQL
func Connect(config *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Database.Host,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.Port,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
