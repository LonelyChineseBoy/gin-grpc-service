package main

import (
	"fmt"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"time"
	"user-srv/handler"
	"user-srv/initialize"
	"user-srv/proto"
)

func init() {
	initialize.InitLogger()
	initialize.InitNacosConfig()
	initialize.InitGormConfig()
}

func randPhone() string {
	rand.Seed(time.Now().UnixNano())

	// 定义运营商代码列表，用于随机选择一个代码
	netIDs := []string{"130", "131", "132", "133", "134", "135", "136", "137", "138", "139", "147", "150", "151", "152", "153", "155", "156", "157", "158", "159", "172", "178", "182", "183", "184", "187", "188", "198"}

	netID := netIDs[rand.Intn(len(netIDs))]
	userNumber := rand.Intn(89999999) + 10000000
	phoneNumber := fmt.Sprintf("%s%d", netID, userNumber)
	return phoneNumber
}
func main() {
	//err := global.DB.AutoMigrate(&model.User{}, &model.UserAddress{})
	//if err != nil {
	//	panic(err)
	//}
	//for i := 0; i < 100; i++ {
	//	phone := randPhone()
	//	password, _ := handler.HashPassword(phone)
	//	global.DB.Create(&model.User{
	//		UserName: "wyy" + phone,
	//		NickName: "nick" + phone,
	//		Password: password,
	//		Mobile:   phone,
	//		Email:    phone + "@qq.com",
	//	})
	//}
	listen, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		panic("listen failed ")
	}
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	err = server.Serve(listen)
	if err != nil {
		panic("grpc server failed")
	}
}
