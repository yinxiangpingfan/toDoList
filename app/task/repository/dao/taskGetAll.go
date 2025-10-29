package dao

import "toDoList/app/task/repository/model"

// 获取用户所有的task
func (db *Dber) GetAllTask(id uint) ([]model.Task, int, error) {
	var tasks []model.Task
	err := db.DB.Where("user_id = ?", id).Find(&tasks).Error
	if err != nil {
		return nil, 1, err
	}
	return tasks, 0, nil
}
