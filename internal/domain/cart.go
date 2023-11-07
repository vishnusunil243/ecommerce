package domain

type Carts struct {
	Id       uint `gorm:"primaryKey;unique;not null"`
	User_id  uint
	Users    Users `gorm:"foreignKey:User_id"`
	CouponId uint
	SubTotal int
	Total    int
}

type CartItem struct {
	Id             uint `gorm:"primaryKey;unique;not null"`
	Carts_id       uint
	Carts          Carts `gorm:"foreignKey:Carts_id"`
	ProductItem_id uint
	ProductItem    ProductItem `gorm:"foreignKey:ProductItem_id"`
	Quantity       int
}
