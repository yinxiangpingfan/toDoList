package dao

import (
	"encoding/json"
	"toDoList/global/logger"

	"github.com/streadway/amqp"
)

// TaskMessage 任务消息结构
type TaskMessage struct {
	Action  string `json:"action"` // add, update, delete
	UserID  uint   `json:"user_id"`
	TaskID  uint   `json:"task_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// PublishTaskMessage 发布任务消息到队列
func PublishTaskMessage(msg TaskMessage) error {
	// 将消息序列化为JSON
	body, err := json.Marshal(msg)
	if err != nil {
		logger.Logger.Errorf("序列化任务消息失败: %v", err)
		return err
	}

	// 发布消息到队列
	err = MqChannel.Publish(
		"",             // exchange（使用默认交换机）
		TaskQueue.Name, // routing key（队列名）
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent, // 消息持久化
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		logger.Logger.Errorf("发布任务消息失败: %v", err)
		return err
	}
	return nil
}
