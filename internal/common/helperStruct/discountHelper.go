package helperStruct

import (
	"time"
)

type Discount struct {
	DiscountPercent   float64
	MaxDiscountAmount int
	MinPurchaseAmount int
	BrandId           uint
	ExpiryDate        time.Time
}
