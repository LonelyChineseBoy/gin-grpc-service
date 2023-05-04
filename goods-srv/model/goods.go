package model

import "gorm.io/gorm"

type Banner struct {
	ID     uint   `gorm:"primarykey"`
	Title  string `gorm:"type:varchar(50);not null;comment:'banner名称'"`
	Url    string `gorm:"not null;comment:'banner图片链接'"`
	Desc   string `gorm:"type:varchar(100);comment:'banner介绍'"`
	IsShow bool   `gorm:"not null;comment:'banner是否被展示'"`
}

type ShopCar struct {
	gorm.Model
	GoodsId    uint32 `gorm:"not null;comment:'加购商品id'"`
	GoodsTotal uint32 `gorm:"not null;comment:'加购商品数量'"`
	UserId     uint   `gorm:"not null;index;comment:'用户id'"`
}
