package response

import (
	"time"
)

type OrderResponse struct {
	OrderDate     time.Time
	PaymentTypeId uint
	PaymentType   string
	Address       `gorm:"embedded" json:"ShippingAddress,omitempty"`
	OrderTotal    int
	OrderStatusID uint
	OrderStatus   string
	PaymentStatus string
	CouponCode    string `json:"coupon,omitempty"`
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
