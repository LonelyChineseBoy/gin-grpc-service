package initialize

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"goods-srv/global"
	"goods-srv/handler"
	"goods-srv/logger"
	"os"
)

func InitLogger() {
	encoder := logger.GetEncoder()
	core := zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel)
	global.Logger = zap.New(core)
	zap.ReplaceGlobals(global.Logger)
}

func InitNacosConfig() {
	handler.ReadConfigByYamlFile("config/nacos.yaml", &global.NacosConfig)
}
