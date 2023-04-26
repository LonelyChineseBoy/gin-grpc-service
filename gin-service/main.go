package main

import (
	"fmt"
	"gin-test/api"
	"gin-test/global"
	"gin-test/handler"
	"gin-test/initialize"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func init() {
	initialize.InitGinLogger()
	initialize.InitGinRouter()
	initialize.InitNacosConfig()
	initialize.InitSentinelConfig()
	api.InitApi()
}
func main() {
	client := handler.NewNacosConfigClient()
	config, _ := client.GetConfig(vo.ConfigParam{
		DataId: global.SentinelConfigParam.DataId,
		Group:  global.SentinelConfigParam.Group,
	})
	fmt.Printf("在nacos中的sentinel配置信息：%s", config)
	global.Router.Run(":8000")
}
