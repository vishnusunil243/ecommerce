package response

import "time"

type OrderResponse struct {
	OrderDate     time.Time
	PaymentTypeId uint
	PaymentType   string
	Address       `gorm:"embedded" json:"ShippingAddress,omitempty"`
	OrderTotal    int
	OrderStatusID uint
	OrderStatus   string
	ProductItemId uint   `json:"ProductItemId,omitempty"`
	ProductName   string `json:"ProductName,omitempty"`
	PaymentStatus string
	CouponCode    string `json:"coupon,omitempty"`
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
