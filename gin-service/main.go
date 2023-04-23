package main

import (
	"fmt"
	"gin-test/api"
	"gin-test/global"
	"gin-test/initialize"
)

func init() {
	initialize.InitGinLogger()
	initialize.InitGinRouter()
	initialize.InitNacos()
	fmt.Println(global.NacosConfig)
	api.InitApi()
}
func main() {
	global.Router.Run(":8000")
}