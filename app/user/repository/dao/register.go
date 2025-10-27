package dao

import (
	"crypto/md5"
	"encoding/hex"
	"toDoList/app/user/repository/model"

	"github.com/go-sql-driver/mysql"
)

func (d *Dber) Register(name string, password string, salt string) (int, error) {
	user := model.User{
		Username: name,
		Password: GenMD5WithSalt(password, salt),
		Salt:     salt,
	}
	result := d.Create(&user)
	if result.Error != nil {
		if mysqlErr, ok := result.Error.(*mysql.MySQLError); ok {
			switch mysqlErr.Number {
			case 1062: // MySQL中表示重复条目的代码
				return 1, mysqlErr
			default:
				return 2, mysqlErr
			}
		} else {
			return 3, result.Error
		}
	}
	return 0, nil
}
func GenMD5WithSalt(passwd, salt string) string {
	s := passwd + "::" + salt
	md5Hash := md5.New()
	md5Hash.Write([]byte(s))
	// 转16进制
	return hex.EncodeToString(md5Hash.Sum(nil))
}
