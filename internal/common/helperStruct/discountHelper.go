package helperStruct

import (
	"time"
)

type Discount struct {
	DiscountPercent float64
	BrandId         uint
	ExpiryDate      time.Time
}
