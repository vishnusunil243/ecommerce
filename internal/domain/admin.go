package domain

import "time"

type Admins struct {
	ID           uint   `gorm:"primaryKey;unique;not null"`
	Name         string `gorm:"not null"`
	Email        string `gorm:"not null;unique"`
	Password     string `gorm:"not null"`
	IsSuperAdmin bool   `gorm:"default:false"`
	IsBlocked    bool   `gorm:"default:false"`
	CreatedAt    time.Time
	UpdateAt     time.Time
}
