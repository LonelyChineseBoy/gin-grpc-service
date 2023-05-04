package main

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"go.uber.org/zap"
	"goods-srv/handler"
	"goods-srv/initialize"
	"goods-srv/proto"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	initialize.InitLogger()
	initialize.InitNacosConfig()
	initialize.InitGormConfig()
}
func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	ip := "127.0.0.1"
	port, _ := handler.GetAvailablePort()
	ClusterName := "goodsServerCluster"
	ServiceName := "goodsServer-" + ip
	GroupName := "goods-group"
	address := fmt.Sprintf("%s:%d", ip, port)
	namingClient := handler.NewNacosNamingClient()
	_, registerErr := namingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		ClusterName: ClusterName,
		ServiceName: ServiceName,
		GroupName:   GroupName,
		Ephemeral:   true,
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
			return
		}
	}
}
