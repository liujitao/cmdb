package common

import (
	"golang.org/x/crypto/bcrypt"
)

// 用户密码明文加密
func SetPassword(password string) string {
	bytePassword := []byte(password)
	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	return string(passwordHash)
}

// 校验用户密码
