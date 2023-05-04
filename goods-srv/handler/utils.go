package handler

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"math/rand"
	"net"
	"strconv"
	"time"
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

// 随机获取可用端口号
func GetAvailablePort() (int, error) {
	rand.Seed(time.Now().UnixNano())
	minPort := 1024
	maxPort := 65535
	port := rand.Intn(maxPort-minPort) + minPort
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	actualPort := listener.Addr().(*net.TCPAddr).Port
	return actualPort, nil
}
