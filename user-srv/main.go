package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ip := "127.0.0.1"
	port, _ := handler.GetAvailablePort()
	ClusterName := "userServerCluster"
	ServiceName := "userServer-" + ip
	GroupName := "user-group"
	address := fmt.Sprintf("%s:%d", ip, port)
	namingClient := handler.NewNacosNamingClient()
	_, registerErr := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Metadata:    nil,
		ClusterName: ClusterName,
		ServiceName: ServiceName,
		GroupName:   GroupName,
		Ephemeral:   false,
	})
	if registerErr != nil {
		zap.S().Panicf("注册用户服务实例失败:%v", registerErr.Error())
	}
	listen, err := net.Listen("tcp", address)
	if err != nil {
		panic("listen failed ")
	}
	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	go func() {
		err = server.Serve(listen)
		if err != nil {
			panic("grpc server failed")
		}
	}()
	for {
		select {
		case <-interrupt:
			server.Stop()
			<-time.After(time.Second * 5)
			server.GracefulStop()
			listen.Close()
			_, deRegisterErr := namingClient.DeregisterInstance(vo.DeregisterInstanceParam{
				Ip:          ip,
				Port:        uint64(port),
				Cluster:     ClusterName,
				ServiceName: ServiceName,
				GroupName:   GroupName,
				Ephemeral:   false,
			})
			if deRegisterErr != nil {
				zap.S().Errorf("注销服务实例失败%v", deRegisterErr)
			}
			return
		}
	}
}
