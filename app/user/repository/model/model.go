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

type User struct {
	gorm.Model
	ID       uint   `gorm:"primaryKey;auto_increment"`
	Username string `gorm:"unique"`
	Password string
	Salt     string
}
