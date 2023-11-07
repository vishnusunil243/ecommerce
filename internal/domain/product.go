package domain

import "time"

type Category struct {
	Id           uint   `gorm:"primaryKey;unique;not null"`
	CategoryName string `gorm:"unique;not null"`
	Created_at   time.Time
	Updated_at   time.Time
}

type Product struct {
	Id          uint   `gorm:"primaryKey;unique;not null"`
	ProductName string `gorm:"unique;not null"`
	Description string
	Brand       string
	Category_id uint
	Category    Category `gorm:"foreignKey:Category_id"`
	Created_at  time.Time
	Updated_at  time.Time
}

type ProductItem struct {
	Id           uint `gorm:"primaryKey;unique;not null"`
	Product_id   uint
	Product      Product `gorm:"foreignKey:Product_id"`
	Sku          string  `gorm:"not null"`
	Qty_in_stock int
	Color        string
	Ram          int
	Battery      int
	Screen_size  float64
	Storage      int
	Camera       int
	Price        int
	Imag         string
	Created_at   time.Time
	Updated_at   time.Time
}

type Images struct {
	Id            uint `gorm:"primaryKey;unique;not null"`
	ProductItemId uint
	ProductItem   ProductItem `gorm:"foreignKey:ProductItemId"`
	FileName      string
}
