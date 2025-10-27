package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"toDoList/app/user/repository/dao"
	"toDoList/app/user/service/pb"
	"toDoList/global/logger"
)

func (u *UserSrv) Register(ctx context.Context, req *pb.RegisterRequest, resp *pb.RegisterResponse) error {
	// 生成随机 salt
	salt, err := generateSalt()
	if err != nil {
		resp.Code = 500
		logger.Logger.Errorf("生成salt发生错误：" + err.Error())
		resp.Msg = "注册失败"
		return err
	}

	// 调用 dao 层注册
	dber := dao.NewUserDBer()
	code, err := dber.Register(req.Username, req.Password, salt)

	if err != nil {
		switch code {
		case 1:
			resp.Code = 400
			resp.Msg = "用户名已存在"
		case 2:
			resp.Code = 500
			resp.Msg = "注册失败"
			logger.Logger.Errorf("注册时发生错误：" + err.Error())
		case 3:
			resp.Code = 500
			logger.Logger.Errorf("注册时发生错误：" + err.Error())
			resp.Msg = "注册失败"
		}
		return err
	}

	resp.Code = 200
	resp.Msg = "注册成功"
	return nil
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}
