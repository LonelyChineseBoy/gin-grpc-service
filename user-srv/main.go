package main

import (
	"golang.org/x/crypto/bcrypt"
	"user-srv/initialize"
)

func init() {
	initialize.InitLogger()
	initialize.InitNacosConfig()
	initialize.InitGormConfig()
}

// 加密密码
func HashPassword(password string) (string, error) {
	// 生成密码哈希值，参数为哈希强度，范围在4~31之间，建议设置为10~14
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 比较密码哈希值是否匹配
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func main() {
	//err := global.DB.AutoMigrate(&model.User{}, &model.UserAddress{})
	//if err != nil {
	//	panic(err)
	//}
}
