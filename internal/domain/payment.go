package domain

import "time"

type PaymentDetails struct {
	OrdersId        int
	OrderTotal      int
	PaymentTypeId   int
	PaymentStatusId int
	UpdatedAt       time.Time
	PaymentType     PaymentType   `gorm:"foreignKey:PaymentTypeId"`
	PaymentStatus   PaymentStatus `gorm:"foreignKey:PaymentStatusId"`
}
type PaymentStatus struct {
	Id     int
	Status string
}
type PaymentType struct {
	Id   uint   `gorm:"primaryKey;unique;not null"`
	Type string `gorm:"unique;not null"`
}
