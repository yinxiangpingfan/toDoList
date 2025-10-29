package main

import (
	"fmt"
	"strconv"
	"toDoList/app/task/repository/dao"
	"toDoList/app/task/service"
	"toDoList/app/task/service/pb"
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
	// 创建ectd
	etcdReg := etcd.NewRegistry(
		registry.Addrs(config.Conf.Service.EtcdHost + ":" + strconv.Itoa(config.Conf.Service.EtcdPort)),
	)
	// 创建微服务实例
	microService := micro.NewService(
		micro.Name("rpcUserService"),
		micro.Address(config.Conf.Service.TaskGrpcHost+":"+strconv.Itoa(config.Conf.Service.TaskGrpcPort)),
		micro.Registry(etcdReg),
	)
	// 初始化服务
	microService.Init()

	// 注册TaskService服务handler
	err := pb.RegisterTaskServiceHandler(microService.Server(), new(service.TaskSrv))
	if err != nil {
		logger.Logger.Panicf("注册TaskService失败: %v", err)
	}

	fmt.Println("TaskService启动成功")

	// 启动服务
	if err := microService.Run(); err != nil {
		logger.Logger.Panicf("启动TaskService失败: %v", err)
	}
}
