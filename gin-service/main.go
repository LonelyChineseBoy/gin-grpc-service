package main

import (
	"gin-test/api"
	"gin-test/global"
	"gin-test/initialize"
)

func init() {
	initialize.InitGinLogger()
	initialize.InitGinRouter()
	api.InitApi()
}
func main() {
	global.Router.Run(":8000")
}
