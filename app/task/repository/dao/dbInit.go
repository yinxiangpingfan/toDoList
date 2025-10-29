package dao

import (
	"fmt"
	"toDoList/app/task/repository/model"
	"toDoList/global/config"
	"toDoList/global/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//初始化task的数据库

type Dber struct {
	*gorm.DB
}

func NewUserDBer() *Dber {
	return &Dber{model.DB}
}

func DatabaseInit() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.Mysql.User,
		config.Conf.Mysql.Password,
		config.Conf.Mysql.Host,
		config.Conf.Mysql.Port,
		config.Conf.Mysql.Database,
	)
	var err error
	model.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.Panicf("数据库连接失败: %v", err)
	}
	model.DB.AutoMigrate(&model.Task{})
}
