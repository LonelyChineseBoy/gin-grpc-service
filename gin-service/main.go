package main

import (
	"gin-service/api"
	"gin-service/global"
	"gin-service/handler"
	"gin-service/initialize"
)

func init() {
	initialize.InitLogger()
	initialize.InitGinRouter()
	initialize.InitNacosConfig()
	initialize.InitSentinelConfig()
	api.InitApi()
}
func main() {
	handler.ExitProcedure(":8000", global.Router)
}
