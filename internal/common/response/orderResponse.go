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
	OrderStatusID uint
	OrderStatus   string
	PaymentStatus string
	CouponCode    string `json:"coupon,omitempty"`
	SubTotal      int    `json:"SubTotal,omitempty"`
	CouponAmount  int    `json:"couponAmount,omitempty"`
	DiscountPrice int    `json:"discount_price,omitempty"`
	OrderTotal    int
}
type OrderProduct struct {
	ProductItemId uint
	Price         int     `json:"price,omitempty"`
	DiscountPrice float64 `json:"DiscountPrice,omitempty"`
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
