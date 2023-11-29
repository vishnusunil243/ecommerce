package response

import "time"

type Wallet struct {
	Amount int
}
type WalletHistories struct {
	Id                uint `gorm:"primaryKey;unique;not null"`
	UserId            uint
	RecentTransaction string
	Balance           int
	Time              time.Time
}
