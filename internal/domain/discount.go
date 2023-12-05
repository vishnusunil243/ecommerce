package domain

import "time"

type Discount struct {
	Id              uint
	DiscountPercent float64
	BrandId         uint
	Brand           Brand `gorm:"foreignKey:BrandId"`
	ExpiryDate      time.Time
}
type Referrals struct {
	Id         uint
	ReferralId string
	UserId     uint
	Users      Users `gorm:"foreignKey:UserId"`
}
type UserReferrals struct {
	Id         uint
	UserId     uint  `gorm:"unique;not null"`
	Users      Users `gorm:"foreignKey:UserId"`
	ReferredBy uint
}
