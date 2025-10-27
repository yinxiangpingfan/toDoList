package utils

import (
	"time"
	"toDoList/global/config"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Id int
	jwt.StandardClaims
}

// 创建jwt
func CreateJwt(id int) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(time.Duration(config.Conf.Jwt.Expire) * time.Second)
	claims := Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "todoList",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(config.Conf.Jwt.Secret)
	return token, err
}

// 解析jwt
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.Conf.Jwt.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
