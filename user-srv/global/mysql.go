package global

var MysqlConfigParam MysqlConfig

type MysqlConfig struct {
	DataId string `mapstructure:"data-id"`
	Group  string `mapstructure:"group"`
}
