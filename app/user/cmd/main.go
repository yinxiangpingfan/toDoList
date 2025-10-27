package main

import (
	"strconv"
	"toDoList/app/user/repository/dao"
	"toDoList/global/config"
	"toDoList/global/logger"

	"github.com/micro/plugins/v5/registry/etcd"
	"go-micro.dev/v5"
	"go-micro.dev/v5/registry"
)

func main() {
	logger.LoggerInit("../../../logs/logs.log")
	config.ConfInit()
	dao.DatabaseInit()
	etcdReg := etcd.NewRegistry(
		registry.Addrs(config.Conf.Service.EtcdHost + ":" + strconv.Itoa(config.Conf.Service.EtcdPort)),
	)
	// 创建微服务实例
	microService := micro.NewService(
		micro.Name("rpcUserService"),
		micro.Address(config.Conf.Service.UserGrpcHost+":"+strconv.Itoa(config.Conf.Service.UserGrpcPort)),
		micro.Registry(etcdReg),
	)
	// 初始化服务
	microService.Init()
}
