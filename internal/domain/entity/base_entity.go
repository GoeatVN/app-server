package entity

import (
	"gorm.io/gorm"
	"time"
)

// BaseEntity contains common fields for all models
type BaseEntity struct {
	CreatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone" json:"updated_at"`
	CreatedBy string    `json:"created_by"` // Record creator
	UpdatedBy string    `json:"updated_by"` // Record updater
}

// BeforeCreate sets the CreatedAt field to the current time
func (b *BaseEntity) BeforeCreate(tx *gorm.DB) (err error) {
	b.CreatedAt = time.Now()
	return
}

// BeforeUpdate sets the UpdatedAt field to the current time
func (b *BaseEntity) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now()
	return
}
