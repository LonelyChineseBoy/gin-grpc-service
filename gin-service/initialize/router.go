package initialize

import (
	"gin-test/global"
	"gin-test/logger"
	"gin-test/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func InitGinLogger() {
	encoder := logger.GetEncoder()
	core := zapcore.NewCore(encoder, os.Stdout, zapcore.DebugLevel)
	global.Logger = zap.New(core)
	zap.ReplaceGlobals(global.Logger)
}

func InitGinRouter() {
	global.Router = gin.New()
	global.Router.Use(middleware.GinLogger(), middleware.GinRecovery(true), middleware.SentinelWarmUp())
}
