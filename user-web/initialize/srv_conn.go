package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop-api/user-web/global"
	"mxshop-api/user-web/proto"
)

func InitSrvConn2() {
	userConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port, global.ServerConfig.UserSrvConfig.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	fmt.Println(global.ServerConfig)
	if err != nil {
		zap.S().Fatal("【InitSrvConn2】连接【用户服务失败")
	}
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClinet = userSrvClient

}
func InitSrvConn() {
	//从注册中心获取用户的服务信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulConfig.Host, global.ServerConfig.ConsulConfig.Port)

	userSrvHost := ""
	userSrvPort := 0
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(`Service == "user-srv"`)
	fmt.Println(data)
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	if userSrvHost == "" {
		zap.S().Fatal("【InitSrvConn】连接【用户服务失败】")
		return
	}

	Userconn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList]链接用户服务失败", "msg", err.Error())
	}
	//后续的用户服务下线了，改变端口，改变ip都会出现问题，在后面的负载均衡中会有使用解决
	userSrvClient := proto.NewUserClient(Userconn)
	global.UserSrvClinet = userSrvClient
}
