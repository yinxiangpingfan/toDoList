package service

import (
	"context"
	"toDoList/app/task/repository/dao"
	"toDoList/app/task/service/pb"
	"toDoList/global/logger"
)

func (s *TaskSrv) DeleteTask(ctx context.Context, req *pb.DelRequest, resp *pb.DelResponse) error {
	db := dao.NewUserDBer()
	code, err := db.DeleteTask(uint(req.Taskid), uint(req.Userid))
	if err != nil {
		logger.Logger.Errorf("DeleteTask failed: %v", err)
		switch code {
		case 1:
			resp.Code = 400
			resp.Msg = err.Error()
		case 2:
			resp.Code = 500
			resp.Msg = "删除任务失败"
		}
		return err
	}
	resp.Code = 200
	resp.Msg = "删除任务成功"
	return nil
}
