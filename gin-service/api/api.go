package api

import (
	"gin-service/global"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func InitApi() {
	group := global.Router.Group("/v1")
	group.GET("/test", Hello)
}

func Hello(ctx *gin.Context) {
	zap.S().Info("return Hello json", zap.String("funcName", "Hello"))
	ctx.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}
