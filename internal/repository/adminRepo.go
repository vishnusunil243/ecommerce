package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
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
func (c *adminDatabase) ListAllUsers(queryParams helperStruct.QueryParams) ([]response.UserDetails, error) {
	var users []response.UserDetails
	getUsers := `SELECT * FROM users`
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		getUsers = fmt.Sprintf("%s LIMIT %d OFFSET %d", getUsers, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		getUsers = fmt.Sprintf("%s LIMIT 10 OFFSET 0", getUsers)
	}
	err := c.DB.Raw(getUsers).Scan(&users).Error
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

func (c *adminDatabase) ReportUser(UsersId int) (response.UserReport, error) {
	var user response.UserReport

	// Execute the UPDATE query
	result := c.DB.Exec("UPDATE users SET report_count=report_count+1 WHERE id=? RETURNING report_count, name", UsersId)

	// Check for errors
	if result.Error != nil {
		return user, result.Error
	}

	// Check the number of rows affected
	rowsAffected := result.RowsAffected
	if rowsAffected == 0 {
		return user, fmt.Errorf("no rows were updated")
	}

	// Manually query for the updated values
	err := c.DB.Raw("SELECT report_count, name FROM users WHERE id = ?", UsersId).Scan(&user).Error

	return user, err
}
