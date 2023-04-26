package initialize

import (
	"gin-test/global"
	"gin-test/handler"
	"gin-test/logger"
	"gin-test/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitLogger() {
	encoder := logger.GetEncoder()
	core := zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel)
	global.Logger = zap.New(core)
	zap.ReplaceGlobals(global.Logger)
}

func InitGinRouter() {
	global.Router = gin.New()
	global.Router.Use(cors.Default())
	global.Router.Use(middleware.GinLogger())
	global.Router.Use(middleware.GinRecovery(true))
	global.Router.Use(middleware.SentinelWarmUp())
}

func InitNacosConfig() {
	handler.ReadConfigByYamlFile("config/nacos.yaml", &global.NacosConfig)
}

func InitSentinelConfig() {
	handler.ReadConfigByYamlFile("config/sentinel.yaml", &global.SentinelConfigParam)
}
