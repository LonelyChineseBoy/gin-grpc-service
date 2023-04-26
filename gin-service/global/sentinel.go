package global

var SentinelConfigParam SentinelConfig

type SentinelConfig struct {
	DataId string `mapstructure:"data-id"`
	Group  string `mapstructure:"group"`
}
