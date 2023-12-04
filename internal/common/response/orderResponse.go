package response

import (
	"time"
)

type OrderResponse struct {
	Id            uint
	OrderDate     time.Time
	PaymentTypeId uint
	PaymentType   string
	Address       `gorm:"embedded" json:"ShippingAddress,omitempty"`
	OrderTotal    int
	OrderStatusID uint
	OrderStatus   string
	PaymentStatus string
	CouponCode    string `json:"coupon,omitempty"`
	CouponAmount  int    `json:"couponAmount,omitempty"`
	SubTotal      int    `json:"SubTotal,omitempty"`
}
type OrderProduct struct {
	ProductItemId uint
	ProductName   string
	Quantity      int
}
type ResponseOrder struct {
	OrderResponse OrderResponse
	OrderProducts []OrderProduct
}
type ReturnOrder struct {
	OrderDate     time.Time
	OrderTotal    int
	OrderStatusID uint
	RefundStatus  string
	OrderStatus   string
}
type AdminOrder struct {
	OrderId       uint
	PaymentTypeId uint
	PaymentType   string
	OrderStatus   string
	PaymentStatus string
}
type OrderStatus struct {
	Id     uint
	Status string
}
