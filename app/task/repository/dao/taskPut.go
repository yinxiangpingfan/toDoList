package dao

import (
	"errors"
	"toDoList/app/task/repository/model"
)

//更新task的任务

func (db *Dber) PutTask(taskId uint, userId uint, title string, content string) (int, error) {
	task := model.Task{}
	//根据taskId查找task
	err := db.Where("id = ? AND user_id = ?", taskId, userId).First(&task).Error
	if err != nil {
		return 1, errors.New("任务不存在或无权限删除")
	}
	task.Title = title
	task.Content = content
	err = db.Save(&task).Error
	if err != nil {
		return 2, err
	}
	return 0, nil
}
