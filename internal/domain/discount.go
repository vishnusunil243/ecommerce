package domain

import "time"

type Discount struct {
	Id                uint
	DiscountPercent   float64
	BrandId           uint
	Brand             Brand `gorm:"foreignKey:BrandId"`
	MaxDiscountAmount int
	MinPurchaseAmount int
	ExpiryDate        time.Time
}
