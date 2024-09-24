package entity

import (
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

func Init(database *gorm.DB) {
	db = database
}

type BaseEntity struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	CreatedBy string    `gorm:"size:255;not null " json:"createdBy"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedBy string    `gorm:"size:255;not null " json:"updatedBy"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
	// DeletedAt *time.Time `sql:"index"`
}

//	func (base *BaseEntity) BeforeCreate(db *gorm.DB) error {
//		base.Id = uuid.NewString()
//		return nil
//	}
func (base *BaseEntity) Prepare() {
	// currentDateTime lấy time theo giờ UTC+7
	currentDateTime := time.Now().UTC().Add(time.Hour * 7)
	base.CreatedAt = currentDateTime
	base.UpdatedAt = currentDateTime
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		if pageSize > 100 {
			pageSize = 100
		}
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
