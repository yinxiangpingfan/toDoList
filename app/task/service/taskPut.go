package service

import (
	"context"
	"toDoList/app/task/repository/dao"
	"toDoList/app/task/service/pb"
	"toDoList/global/logger"
)

// UpdateTask 更新任务（注意：方法名必须与proto定义一致）
func (s *TaskSrv) UpdateTask(ctx context.Context, req *pb.UpdateRequest, resp *pb.UpdateResponse) error {
	db := dao.NewUserDBer()
	code, err := db.PutTask(uint(req.Taskid), uint(req.Userid), req.Title, req.Content)
	if err != nil {
		logger.Logger.Errorf("UpdateTask error: %v", err)
		switch code {
		case 1:
			resp.Code = 400
			resp.Msg = "更新任务失败"
		case 2:
			resp.Code = 500
			resp.Msg = "更新任务失败"
		}
		return err
	}
	resp.Code = 200
	resp.Msg = "更新任务成功"
	return nil
}
