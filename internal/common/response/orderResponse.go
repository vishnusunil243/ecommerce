package response

import "time"

type OrderResponse struct {
	OrderDate     time.Time
	PaymentTypeId uint
	PaymentType   string
	Address       `gorm:"embedded" json:"ShippingAddress"`
	OrderTotal    int
	OrderStatusID uint
	OrderStatus   string
	CouponCode    string
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
