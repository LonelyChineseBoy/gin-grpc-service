package global

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	Router *gin.Engine
	Logger *zap.Logger
)
