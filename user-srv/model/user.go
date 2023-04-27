package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	UserName  string    `gorm:"type:varchar(50);not null;unique;comment:'用户名'"`
	NickName  string    `grom:"type:varchar(50);comment:'用户昵称'"`
	Password  string    `gorm:"type:varchar(1024);not null;comment:'用户密码'"`
	Mobile    string    `gorm:"type:varchar(15);not null;unique;comment:'用户手机号'"`
	Email     string    `gorm:"type:varchar(50);comment:'用户邮箱'"`
	Status    uint8     `gorm:"default:0;comment:'用户状态,0表示正常用户,1表示禁用用户'"`
	UserType  uint8     `gorm:"default:0;comment:'用户类型,0表示普通用户,1表示后台系统管理用户'"`
	Gender    uint8     `gorm:"default:0;comment:'用户性别,0表示男,1表示女'"`
	Birthdaye time.Time `gorm:"comment:'用户生日时间'"`
}

type UserAddress struct {
	gorm.Model
	UserId        int    `gorm:"index;not null;comment:'关联用户id'"`
	AddressDetail string `gorm:"type:varchar(200);not null;comment:'用户地址'"`
}
