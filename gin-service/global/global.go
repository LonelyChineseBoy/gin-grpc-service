package global

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Router *gin.Engine
	Logger *zap.Logger
	DB     *gorm.DB
)
