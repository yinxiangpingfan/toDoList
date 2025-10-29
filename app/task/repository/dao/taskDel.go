package dao

import (
	"errors"
	"toDoList/app/task/repository/model"
)

// 删除task的任务
func (db *Dber) DeleteTask(taskID uint, userID uint) (int, error) {
	// 先查询任务是否存在且属于该用户
	var task model.Task
	err := db.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error
	if err != nil {
		return 1, errors.New("任务不存在或无权限删除")
	}

	// 删除任务
	err = db.DB.Delete(&task).Error
	if err != nil {
		return 2, err
	}

	return 0, nil
}
