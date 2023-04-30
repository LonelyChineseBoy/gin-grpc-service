package handler

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func GetAvailablePort() (int, error) {
	rand.Seed(time.Now().UnixNano())
	minPort := 1024
	maxPort := 65535
	port := rand.Intn(maxPort-minPort) + minPort
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return 0, err
	}
	listener.Close()
	actualPort := listener.Addr().(*net.TCPAddr).Port
	return actualPort, nil
}

// 加密密码
func HashPassword(password string) (string, error) {
	// 生成密码哈希值，参数为哈希强度，范围在4~31之间，建议设置为10~14
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 比较密码哈希值是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
