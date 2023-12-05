package response

import "time"

type Discount struct {
	Id              uint
	DiscountPercent float64
	BrandName       string
	ExpiryDate      time.Time
}
