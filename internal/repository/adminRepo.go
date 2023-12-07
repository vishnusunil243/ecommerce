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
func (c *adminDatabase) ListAllUsers(queryParams helperStruct.QueryParams) ([]response.UserDetails, int, error) {
	var users []response.UserDetails
	getUsers := `SELECT * FROM users`
	var count int
	getTotalCount := fmt.Sprintf("SELECT COUNT(*) FROM (%s)", getUsers)
	err := c.DB.Raw(getTotalCount).Scan(&count).Error
	if err != nil {
		return []response.UserDetails{}, 0, err
	}
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		getUsers = fmt.Sprintf("%s LIMIT %d OFFSET %d", getUsers, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		getUsers = fmt.Sprintf("%s LIMIT 10 OFFSET 0", getUsers)
	}
	err = c.DB.Raw(getUsers).Scan(&users).Error
	return users, count, err
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

// GetDashBoard implements interfaces.AdminRepository.
func (c *adminDatabase) GetDashBoard(dashboard helperStruct.Dashboard) (response.DashBoard, error) {
	tx := c.DB.Begin()
	var newDashboard response.DashBoard
	getDashboard := `SELECT SUM(oi.quantity*oi.price) AS total_revenue,
	                 SUM(oi.quantity) AS total_products_sold,
					  COUNT(DISTINCT o.id) AS total_orders 
					  FROM orders o JOIN order_items oi ON o.id=oi.orders_id`
	getUsers := `SELECT COUNT(id) AS TotalUsers FROM USERS`
	if dashboard.Year != 0 {
		getDashboard = fmt.Sprintf("%s WHERE EXTRACT (YEAR FROM order_date) =%d ", getDashboard, dashboard.Year)
		getUsers = fmt.Sprintf("%s WHERE EXTRACT (YEAR FROM created_at)=%d", getUsers, dashboard.Year)
		if dashboard.Month != 0 {
			getDashboard = fmt.Sprintf("%s AND EXTRACT (MONTH FROM order_date)=%d", getDashboard, dashboard.Month)
			getUsers = fmt.Sprintf("%s AND EXTRACT (MONTH FROM created_at)=%d", getUsers, dashboard.Month)
			if dashboard.Day != 0 {
				getDashboard = fmt.Sprintf("%s AND EXTRACT (DAY FROM order_date)=%d", getDashboard, dashboard.Day)
				getUsers = fmt.Sprintf("%s AND EXTRACT (DAY FROM created_at)=%d", getUsers, dashboard.Day)
			}
		}
	} else if dashboard.StartDate != "" && dashboard.EndDate != "" {
		getDashboard = fmt.Sprintf("%s WHERE o.order_date BETWEEN '%s' AND '%s'", getDashboard, dashboard.StartDate, dashboard.EndDate)
		getUsers = fmt.Sprintf("%s WHERE created_at BETWEEN '%s' AND '%s'", getUsers, dashboard.StartDate, dashboard.EndDate)
	}
	err := tx.Raw(getDashboard).Scan(&newDashboard).Error
	if err != nil {
		tx.Rollback()
		return response.DashBoard{}, fmt.Errorf("error returning dashboard")
	}
	err = tx.Raw(getUsers).Scan(&newDashboard.TotalUsers).Error
	if err != nil {
		tx.Rollback()
		return response.DashBoard{}, err
	}
	return newDashboard, nil
}

// ViewSalesReport implements interfaces.AdminRepository.
func (c *adminDatabase) ViewSalesReport(filter helperStruct.Dashboard) ([]response.SalesReport, error) {
	var salesReports []response.SalesReport
	getSalesReport := `SELECT u.name,pt.type AS payment_type,
	                   o.order_date,o.order_total
					    FROM orders o JOIN users u ON u.id=o.user_id
						 JOIN payment_types pt ON pt.id=o.payment_type_id
						WHERE o.order_status_id=4`
	if filter.Year != 0 {
		getSalesReport = fmt.Sprintf(`%s AND EXTRACT(YEAR FROM o.order_date)=%d`, getSalesReport, filter.Year)
		if filter.Month != 0 {
			getSalesReport = fmt.Sprintf(`%s AND EXTRACT(MONTH FROM o.order_date)=%d `, getSalesReport, filter.Month)
		}
		if filter.Day != 0 {
			getSalesReport = fmt.Sprintf(`%s AND EXTRACT(DAY FROM o.order_date)=%d`, getSalesReport, filter.Day)
		}
	} else if filter.StartDate != "" && filter.EndDate != "" {
		getSalesReport = fmt.Sprintf(`%s AND o.order_date BETWEEN '%s' AND '%s'`, getSalesReport, filter.StartDate, filter.EndDate)
	} else {
		getSalesReport = fmt.Sprintf("%s ORDER BY o.order_date DESC", getSalesReport)
	}
	err := c.DB.Raw(getSalesReport).Scan(&salesReports).Error
	return salesReports, err
}
