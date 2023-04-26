package main

import (
	"gin-test/api"
	"gin-test/global"
	"gin-test/handler"
	"gin-test/initialize"
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
