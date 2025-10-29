package model

import (
	"gorm.io/gorm"
)

var DB *gorm.DB

type Dber struct {
	*gorm.DB
}

var DBer = Dber{
	DB,
}

type Task struct {
	gorm.Model
	ID      uint   `gorm:"primaryKey;autoIncrement"`
	UserID  uint   `gorm:"not null;index"`
	Title   string `gorm:"type:varchar(200);not null"`
	Content string `gorm:"type:text"`
}
