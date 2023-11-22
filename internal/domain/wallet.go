package domain

type Wallet struct {
	Id     uint `gorm:"primaryKey;unique;not null"`
	UserId uint
	Users  Users `gorm:"foreignKey:UserId"`
	Amount int
}
