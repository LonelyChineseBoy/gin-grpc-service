package main

import (
	"user-srv/global"
	"user-srv/initialize"
	"user-srv/model"
)

func init() {
	initialize.InitLogger()
	initialize.InitNacosConfig()
	initialize.InitGormConfig()
}
func main() {
	err := global.DB.AutoMigrate(&model.User{}, &model.UserAddress{})
	if err != nil {
		panic(err)
	}
}
