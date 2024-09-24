package repository

import (
	"food-app/domain/entity"
)

type SysConfigRepository interface {
	SaveUser(config *entity.SysConfig) (*entity.SysConfig, map[string]string)
	GetSysConfig(string) (*entity.SysConfig, error)
	GetSysConfigs() ([]entity.SysConfig, error)
}
