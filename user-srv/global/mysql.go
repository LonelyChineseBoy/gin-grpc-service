package global

import "gorm.io/gorm"

var (
	MysqlConfigParam MysqlConfig
	DB               *gorm.DB
)

type MysqlConfig struct {
	DataId string `mapstructure:"data-id"`
	Group  string `mapstructure:"group"`
}
