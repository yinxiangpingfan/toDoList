package dao

import (
	"toDoList/app/task/repository/model"
)

// 增加task的任务
func (db *Dber) AddTask(userID uint, title string, content string) (uint, error) {
	task := model.Task{
		UserID:  userID,
		Title:   title,
		Content: content,
	}

	err := db.DB.Create(&task).Error
	if err != nil {
		return 1, err
	}

	return 0, nil
}
