package dao

// AddTaskAsync 异步创建任务（通过RabbitMQ）
func (db *Dber) AddTaskAsync(userID uint, title string, content string) error {
	// 构造任务消息
	msg := TaskMessage{
		Action:  "add",
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	// 发布到RabbitMQ队列
	return PublishTaskMessage(msg)
}

// UpdateTaskAsync 异步更新任务（通过RabbitMQ）
func (db *Dber) UpdateTaskAsync(taskID uint, userID uint, title string, content string) error {
	msg := TaskMessage{
		Action:  "update",
		UserID:  userID,
		TaskID:  taskID,
		Title:   title,
		Content: content,
	}

	return PublishTaskMessage(msg)
}

// DeleteTaskAsync 异步删除任务（通过RabbitMQ）
func (db *Dber) DeleteTaskAsync(taskID uint, userID uint) error {
	msg := TaskMessage{
		Action: "delete",
		UserID: userID,
		TaskID: taskID,
	}

	return PublishTaskMessage(msg)
}
