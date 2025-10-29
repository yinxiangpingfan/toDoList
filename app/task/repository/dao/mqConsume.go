package dao

import (
	"encoding/json"
	"toDoList/global/logger"
)

// StartTaskConsumer 启动任务消费者
func StartTaskConsumer() {
	msgs, err := MqChannel.Consume(
		TaskQueue.Name, // 队列名
		"",             // 消费者名字（自动生成）
		false,          // 手动确认消息
		false,          // 非独占
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		logger.Logger.Panicf("启动任务消费者失败: %v", err)
	}

	// 创建一个channel用于优雅关闭
	forever := make(chan bool)

	// 启动goroutine处理消息
	go func() {
		for msg := range msgs {
			processTaskMessage(msg.Body)
			// 手动确认消息已处理
			msg.Ack(false)
		}
	}()

	<-forever
}

// processTaskMessage 处理任务消息
func processTaskMessage(body []byte) {
	var taskMsg TaskMessage
	err := json.Unmarshal(body, &taskMsg)
	if err != nil {
		logger.Logger.Errorf("解析任务消息失败: %v", err)
		return
	}

	db := NewUserDBer()

	switch taskMsg.Action {
	case "add":
		// 创建任务
		_, err := db.AddTask(taskMsg.UserID, taskMsg.Title, taskMsg.Content)
		if err != nil {
			logger.Logger.Errorf("创建任务失败: %v", err)
		}

	case "update":
		// 更新任务
		_, err := db.PutTask(taskMsg.TaskID, taskMsg.UserID, taskMsg.Title, taskMsg.Content)
		if err != nil {
			logger.Logger.Errorf("更新任务失败: %v", err)
		}

	case "delete":
		// 删除任务
		_, err := db.DeleteTask(taskMsg.TaskID, taskMsg.UserID)
		if err != nil {
			logger.Logger.Errorf("删除任务失败: %v", err)
		}

	default:
		logger.Logger.Warnf("未知的任务操作: %s", taskMsg.Action)
	}
}
