package response

import "time"

type Discount struct {
	Id                uint
	MinPurchaseAmount int
	MaxDiscountAmount int
	DiscountPercent   float64
	BrandName         string
	ExpiryDate        time.Time
}
