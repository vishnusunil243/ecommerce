package domain

type Wishlist struct {
	Id             uint `gorm:"primaryKey;unique;not null"`
	User_id        uint
	Users          Users `gorm:"foreignKey:User_id"`
	ProductItem_id uint
	ProductItem    ProductItem `gorm:"foreignKey:ProductItem_id"`
}
