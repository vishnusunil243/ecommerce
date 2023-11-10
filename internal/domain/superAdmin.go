package domain

type SuperAdmin struct {
	Id       uint `gorm:"primaryKey"`
	Name     string
	Email    string
	Password string
}
