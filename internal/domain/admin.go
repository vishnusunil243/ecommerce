package domain

import "time"

type Admins struct {
	ID        uint   `gorm:"primaryKey;unique;not null"`
	Name      string `gorm:"not null"`
	Email     string `gorm:"not null;unique"`
	Password  string `gorm:"not null"`
	IsBlocked bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdateAt  time.Time
}
type AdminInfo struct {
	ID         uint `gorm:"primaryKey"`
	AdminId    uint
	BlockedAt  time.Time
	BlockUntil time.Time
}
