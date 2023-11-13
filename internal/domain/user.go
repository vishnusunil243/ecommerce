package domain

import "time"

type Users struct {
	ID          uint   `gorm:"primarykey;unique;notnull"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required,email" gorm:"unique;not null"`
	Mobile      string `json:"mobile" binding:"required,eq=10" gorm:"unique;not null"`
	Password    string `json:"password" gorm:"not null"`
	IsBlocked   bool   `gorm:"default:false"`
	ReportCount int
	CreatedAt   time.Time
}
type Address struct {
	ID           uint `gorm:"primaryKey;unique;not null"`
	UsersID      uint
	Users        Users  `gorm:"foreignKey:UsersID"`
	House_number string `json:"house_number" binding:"required"`
	Street       string `json:"street" binding:"required"`
	City         string `json:"city " binding:"required"`
	District     string `json:"district " binding:"required"`
	Landmark     string `json:"landmark" binding:"required"`
	Pincode      int    `json:"pincode " binding:"required"`
	IsDefault    bool   `gorm:"default:false"`
}
type UserInfo struct {
	ID                uint `gorm:"primaryKey"`
	UsersID           uint
	Users             Users `gorm:"foreignKey:UsersID"`
	BlockedAt         time.Time
	BlockUntil        time.Time
	BlockedBy         uint
	ReportCount       uint
	ReasonForBlocking string
}
type ReportInfo struct {
	ID                 uint `gorm:"primaryKey"`
	UsersId            uint
	Users              Users `gorm:"foreignKey:UsersId"`
	ReportCount        uint
	ReasonForReporting string
	ReportedBy         uint
}
