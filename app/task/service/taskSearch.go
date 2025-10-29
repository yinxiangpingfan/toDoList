package service

import (
	"context"
	"toDoList/app/task/repository/dao"
	"toDoList/app/task/service/pb"
	"toDoList/global/logger"
)

// GetAllTasks 获取用户的所有任务（注意：方法名必须与proto定义一致）
func (t *TaskSrv) GetAllTasks(ctx context.Context, req *pb.GetAllRequest, resp *pb.GetAllResponse) error {
	db := dao.NewUserDBer()
	tasks, code, err := db.GetAllTask(uint(req.Userid))
	if err != nil {
		logger.Logger.Errorf("GetAllTasks: %v", err)
		switch code {
		case 1:
			resp.Tasks = nil
			resp.Code = 400
			resp.Msg = "查询失败"
			return err
		}
		return err
	}

	// 将 []model.Task 转换为 []*pb.TaskInfo
	taskInfos := make([]*pb.TaskInfo, 0, len(tasks))
	for _, task := range tasks {
		taskInfos = append(taskInfos, &pb.TaskInfo{
			Id:        uint64(task.ID),
			Title:     task.Title,
			Content:   task.Content,
			CreatedAt: task.CreatedAt.Unix(),
			UpdatedAt: task.UpdatedAt.Unix(),
		})
	}

	resp.Tasks = taskInfos
	resp.Code = 200
	resp.Msg = "查询成功"
	return nil
}
