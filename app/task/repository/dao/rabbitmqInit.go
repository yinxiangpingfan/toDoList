package dao

import (
	"fmt"
	"toDoList/global/config"
	"toDoList/global/logger"

	"github.com/streadway/amqp"
)

var MqConn *amqp.Connection
var MqChannel *amqp.Channel
var TaskQueue amqp.Queue

// RabbitMqInit 初始化RabbitMQ连接和队列
func RabbitMqInit() {
	var err error
	// 连接RabbitMQ
	MqConn, err = amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%d/",
		config.Conf.RabbitMq.Protocol,
		config.Conf.RabbitMq.User,
		config.Conf.RabbitMq.Password,
		config.Conf.RabbitMq.Host,
		config.Conf.RabbitMq.Port))
	if err != nil {
		logger.Logger.Panicf("RabbitMQ连接失败: %s", err.Error())
	}

	// 创建Channel
	MqChannel, err = MqConn.Channel()
	if err != nil {
		logger.Logger.Panicf("RabbitMQ创建Channel失败: %v", err)
	}

	// 声明队列
	TaskQueue, err = MqChannel.QueueDeclare(
		"task_queue", // 队列名字
		true,         // 消息持久化
		false,        // 不自动删除
		false,        // 非独占
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		logger.Logger.Panicf("RabbitMQ声明队列失败: %v", err)
	}

	logger.Logger.Infof("RabbitMQ初始化成功，队列: %s", TaskQueue.Name)
}

// RabbitMqClose 关闭RabbitMQ连接
func RabbitMqClose() {
	if MqChannel != nil {
		MqChannel.Close()
	}
	if MqConn != nil {
		MqConn.Close()
	}
	fmt.Println("RabbitMQ连接已关闭")
}
