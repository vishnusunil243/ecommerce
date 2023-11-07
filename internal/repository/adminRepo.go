package repository

import (
	"gorm.io/gorm"
	"main.go/internal/domain"
	"main.go/internal/repository/interfaces"
)

type adminDatabase struct {
	DB *gorm.DB
}

func NewAdminRepo(DB *gorm.DB) interfaces.AdminRepository {
	return &adminDatabase{DB}

}

// AdminLogin implements interfaces.AdminRepository.
func (c *adminDatabase) AdminLogin(email string) (domain.Admins, error) {
	var adminData domain.Admins
	err := c.DB.Raw("SELECT * FROM admins WHERE email=?", email).Scan(&adminData).Error
	if err != nil {
		return adminData, err
	}
	return adminData, nil
}
