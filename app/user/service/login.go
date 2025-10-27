package service

import (
	"context"
	"toDoList/app/user/repository/dao"
	"toDoList/app/user/service/pb"
	"toDoList/global/logger"
)

func (u *UserSrv) Login(ctx context.Context, req *pb.LoginRequest, res *pb.LoginResponse) error {
	db := dao.NewUserDBer()
	code, err := db.Login(req.Username, req.Password)
	if err != nil {
		switch code {
		case 1:
			res.Code = 400
			res.Msg = err.Error()
		case 2:
			res.Code = 2
			res.Msg = err.Error()
		default:
			res.Code = 3
			logger.Logger.Errorf("登陆时发生错误：" + err.Error())
			res.Msg = err.Error()
		}
		return err
	}
	res.Code = 200
	res.Msg = "登陆成功"
	return nil
}
