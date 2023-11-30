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
type Brand struct {
	Id          uint   `gorm:"primaryKey;unique;not null"`
	Brandname   string `gorm:"unique;not null"`
	Description string
	Category_id uint
	Category    Category `gorm:"foreignKey:Category_id"`
	Created_at  time.Time
	Updated_at  time.Time
}

type ProductItem struct {
	Id                uint `gorm:"primaryKey;unique;not null"`
	Product_id        uint
	Product           Product `gorm:"foreignKey:Product_id"`
	Sku               string  `gorm:"not null"`
	Qty_in_stock      int
	Color             string
	Ram               int
	Battery           int
	Screen_size       float64
	Storage           int
	Camera            int
	Graphic_Processor string
	Price             int
	Created_at        time.Time
	Updated_at        time.Time
}

type Images struct {
	Id            uint `gorm:"primaryKey;unique;not null"`
	ProductItemId uint
	ProductItem   ProductItem `gorm:"foreignKey:ProductItemId"`
	Image         []byte
}
type Image_items struct {
	Id            uint `gorm:"primaryKey;unique;not null"`
	ProductItemId uint
	ProductItem   ProductItem `gorm:"foreignKey:ProductItemId"`
	Image         string
	IsDefault     bool `gorm:"default:false"`
}
