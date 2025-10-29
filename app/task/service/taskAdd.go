package service

import (
	"context"
	"toDoList/app/task/repository/dao"
	"toDoList/app/task/service/pb"
	"toDoList/global/logger"
)

func (s *TaskSrv) AddTask(ctx context.Context, req *pb.AddRequest, resp *pb.AddResponse) error {
	db := dao.NewUserDBer()
	code, error := db.AddTask(uint(req.Id), req.Title, req.Content)
	if error != nil {
		logger.Logger.Error(error.Error())
		switch code {
		case 1:
			resp.Code = 400
			resp.Msg = "增加待办任务失败"
		}
		return error
	}
	resp.Code = 200
	resp.Msg = "增加待办任务成功"
	return nil
}
