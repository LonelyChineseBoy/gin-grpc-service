package handler

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func ReadConfigByYamlFile(filePath string, globalVal interface{}) {
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		zap.S().Errorf("viper read config file failed %s \n", err)
	}
	if err := viper.Unmarshal(globalVal); err != nil {
		zap.S().Errorf("unmarshal conf failed, err:%s \n", err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Info("The configuration file is modified")
		if err := viper.Unmarshal(globalVal); err != nil {
			zap.S().Errorf("unmarshal conf failed, err:%s \n", err)
		}
	})
}
