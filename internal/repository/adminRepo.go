package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/response"
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

// ListAllUsers implements interfaces.AdminRepository.
func (c *adminDatabase) ListAllUsers() ([]response.UserDetails, error) {
	var users []response.UserDetails
	err := c.DB.Raw(`SELECT * FROM users`).Scan(&users).Error
	return users, err
}

// DispalyUser implements interfaces.AdminRepository.
func (c *adminDatabase) DispalyUser(id int) (response.UserDetails, error) {
	var user response.UserDetails
	err := c.DB.Raw(`SELECT * FROM users WHERE id=?`, id).Scan(&user).Error
	if user.Email == "" {
		return user, fmt.Errorf("user not found")
	}
	return user, err
}
