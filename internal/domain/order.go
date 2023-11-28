package domain

import (
	"time"
)

type Orders struct {
	Id              uint `gorm:"primaryKey;unique;not null"`
	UserId          uint
	Users           Users `gorm:"foreignKey:UserId" json:"-"`
	OrderDate       time.Time
	PaymentTypeId   uint
	PaymentType     PaymentType `gorm:"foreignKey:PaymentTypeId" json:"-"`
	ShippingAddress uint
	Address         Address `gorm:"foreignKey:ShippingAddress"`
	OrderTotal      int
	OrderStatusID   uint
	OrderStatus     OrderStatus `gorm:"foreignKey:OrderStatusID" json:"-"`
	PaymentStatusId uint
	PaymentStatus   PaymentStatus `gorm:"foreignKey:PaymentStatusId" json:"-"`
	CouponCode      uint
	Coupon          Coupon `gorm:"foreignKey:CouponCode"`
}

type OrderItem struct {
	Id            uint `gorm:"primaryKey;unique;not null"`
	OrdersId      uint
	Orders        Orders `gorm:"foreignKey:OrdersId" json:"-"`
	ProductItemId uint
	ProductItem   ProductItem `gorm:"foreignKey:ProductItemId" json:"-"`
	Quantity      int
	Price         int
}

type OrderStatus struct {
	Id     uint   `gorm:"primaryKey;unique;not null"`
	Status string `grom:"unique"`
}
