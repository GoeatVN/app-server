package entity

import "time"

type SysConfig struct {
	BaseEntity
	ID          uint64    `gorm:"" json:"id"`
	ConfigGroup string    `gorm:"size:64;not null;" json:"configGroup"`
	ConfigCode  string    `gorm:"size:64;not null;" json:"configCode"`
	ConfigName  string    `gorm:"size:64;not null;unique" json:"configName"`
	ConfigValue string    `gorm:"size:64;not null;" json:"configValue"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"createdAt"`
	CreatedBy   string    `gorm:"size:64;not null;" json:"createdBy"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updatedAt"`
	UpdatedBy   string    `gorm:"size:64;not null;" json:"updatedBy"`
}

func (u *SysConfig) TableName() string {
	return "sys_config"
}
