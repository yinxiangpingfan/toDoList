package dao

import (
	"errors"
	"toDoList/app/user/repository/model"

	"gorm.io/gorm"
)

func (d *Dber) Login(name string, password string) (int, error) {
	user := model.User{}
	err := d.Where("username = ?", name).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 1, errors.New("用户不存在")
		}
		return 3, err
	}
	if GenMD5WithSalt(password, user.Salt) != user.Password {
		return 2, errors.New("密码错误")
	}
	return 0, nil
}
