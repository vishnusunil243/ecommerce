package domain

import "time"

type Coupon struct {
	Id         uint   `gorm:"primaryKey;unique;not null"`
	Name       string `gorm:"unique"`
	Amount     int    `gorm:"CHECK(amount>=0)"`
	Quantity   int    `gorm:"CHECK(quantity>=0)"`
	IsDisabled bool   `gorm:"default:false"`
	CreatedAt  time.Time
}
type UserCoupons struct {
	Id       uint `gorm:"primaryKey;unique;not null"`
	UserId   uint
	Users    Users `gorm:"foreignKey:UserId"`
	CouponId uint
	Coupon   Coupon `gorm:"foreignKey:CouponId"`
	OrderId  uint
	Orders   Orders `gorm:"foreignKey:OrderId"`
}
type UserRewardCoupons struct {
	Id       uint
	UsersId  uint
	Users    Users `gorm:"foreignKey:UsersId"`
	CouponId uint
	Coupon   Coupon `gorm:"foreignKey:CouponId"`
}
