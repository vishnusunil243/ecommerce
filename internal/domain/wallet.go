package domain

type Wallet struct {
	Id     uint `gorm:"primaryKey"`
	UserId uint
	Users  Users `gorm:"foreignKey:UserId"`
	Amount int
}
