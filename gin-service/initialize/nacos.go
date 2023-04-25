package initialize

import (
	"gin-test/global"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func InitNacosConfig() {
	viper.SetConfigFile("config/nacos.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		zap.S().Errorf("viper read config file failed %s \n", err)
	}
	if err := viper.Unmarshal(&global.NacosConfig); err != nil {
		zap.S().Errorf("unmarshal conf failed, err:%s \n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Info("The NACOS configuration file is modified")
		if err := viper.Unmarshal(&global.NacosConfig); err != nil {
			zap.S().Errorf("unmarshal conf failed, err:%s \n", err)
		}
	})
}
