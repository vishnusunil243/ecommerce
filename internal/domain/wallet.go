package domain

import "time"

type Wallet struct {
	Id     uint `gorm:"primaryKey;unique;not null"`
	UserId uint
	Users  Users `gorm:"foreignKey:UserId"`
	Amount int
}
type WalletHistories struct {
	Id                uint `gorm:"primaryKey;unique;not null"`
	UserId            uint
	Users             Users `gorm:"foreignKey:UserId"`
	RecentTransaction string
	Time              time.Time
	Balance           int
}
