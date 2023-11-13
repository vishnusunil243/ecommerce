package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/domain"
	"main.go/internal/repository/interfaces"
)

type superAdminDatabase struct {
	DB *gorm.DB
}

func NewSuperRepo(DB *gorm.DB) interfaces.SuperAdminRepository {
	return &superAdminDatabase{
		DB: DB,
	}
}

// Login implements interfaces.SuperAdminRepository.
func (c *superAdminDatabase) Login(superadmin helperStruct.SuperLoginReq) (domain.SuperAdmin, error) {
	var superAdmin domain.SuperAdmin
	err := c.DB.Raw(`SELECT * FROM super_admins WHERE email=?`, superadmin.Email).Scan(&superAdmin).Error
	return superAdmin, err

}

// CreateAdmin implements interfaces.SuperAdminRepository.
func (c *superAdminDatabase) CreateAdmin(admin helperStruct.CreateAdmin) (response.AdminData, error) {
	var newAdmin response.AdminData
	insertQuery := `INSERT INTO admins(name,email,password,created_at)VALUES($1,$2,$3,NOW()) RETURNING name,email,id`
	err := c.DB.Raw(insertQuery, admin.Name, admin.Email, admin.Password).Scan(&newAdmin).Error
	return newAdmin, err

}

// ListAllAdmins implements interfaces.SuperAdminRepository.
func (c *superAdminDatabase) ListAllAdmins(queryParams helperStruct.QueryParams) ([]response.AdminData, error) {
	var admins []response.AdminData
	getAdmins := `SELECT * FROM admins`
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		getAdmins = fmt.Sprintf("%s LIMIT %d OFFSET %d", getAdmins, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		getAdmins = fmt.Sprintf("%s LIMIT 10 OFFSET 0", getAdmins)
	}
	err := c.DB.Raw(getAdmins).Scan(&admins).Error
	return admins, err
}

// DisplayAdmin implements interfaces.SuperAdminRepository.
func (c *superAdminDatabase) DisplayAdmin(id int) (response.AdminData, error) {
	var admin response.AdminData
	err := c.DB.Raw(`SELECT * FROM admins WHERE id=?`, id).Scan(&admin).Error
	if admin.Email == "" {
		return admin, fmt.Errorf("no admin found with given id")
	}
	return admin, err
}

// BlockAdmin implements interfaces.SuperAdminRepository.
func (c *superAdminDatabase) BlockAdmin(id int) (response.AdminData, error) {
	var exists bool
	var admin response.AdminData
	c.DB.Raw(`select exists(select 1 from admins where id=?)`, id).Scan(&exists)
	if !exists {
		return admin, fmt.Errorf("no  admin found with given id")
	}
	updateQuery := `UPDATE admins SET is_blocked=true WHERE id=? RETURNING id,name,email `
	err := c.DB.Raw(updateQuery, id).Scan(&admin).Error
	if err != nil {
		return admin, err
	}
	err = c.DB.Exec(`
		INSERT INTO admin_infos(admin_id, blocked_at, block_until) 
		(SELECT id, NOW(), NOW() + INTERVAL '5 minutes' FROM admins WHERE id=?) 
	`, admin.Id).Error
	return admin, err

}

// BlockUser implements interfaces.SuperAdminRepository.
func (c *superAdminDatabase) BlockUser(id int) (response.UserData, error) {
	var exists bool
	var userData response.UserData
	c.DB.Raw(`select exists(select 1 from users where id=?)`, id).Scan(&exists)
	if !exists {
		return userData, fmt.Errorf("no  user found with given id")
	}
	err := c.DB.Raw(`UPDATE users SET is_blocked=true WHERE id=? RETURNING id,email,name,mobile`, id).Scan(&userData).Error
	if err != nil {
		return userData, err
	}
	err = c.DB.Exec(`
		INSERT INTO user_infos(users_id, blocked_at, block_until) 
		(SELECT id, NOW(), NOW() + INTERVAL '4 minutes' FROM users WHERE id=?) 
	`, userData.Id).Error
	return userData, err
}

// UnBlockAdminManually implements interfaces.SuperAdminRepository.
func (c *superAdminDatabase) UnBlockAdminManually(id int) (response.AdminData, error) {
	var exists bool
	var admin response.AdminData
	c.DB.Raw(`select exists(select 1 from admins where id=?)`, id).Scan(&exists)
	if !exists {
		return admin, fmt.Errorf("no  admin found with given id")
	}
	updateQuery := `UPDATE admins SET is_blocked=false WHERE id=? RETURNING id,email,name`
	err := c.DB.Raw(updateQuery, id).Scan(&admin).Error
	return admin, err
}

// UnBlockUserManually implements interfaces.SuperAdminRepository.
func (c *superAdminDatabase) UnBlockUserManually(id int) (response.UserData, error) {
	var exists bool
	var userData response.UserData
	c.DB.Raw(`select exists(select 1 from users where id=?)`, id).Scan(&exists)
	if !exists {
		return userData, fmt.Errorf("no  user found with given id")
	}
	updateQuery := `UPDATE users SET is_blocked=false WHERE id=? RETURNING id,email,name,mobile `
	err := c.DB.Raw(updateQuery, id).Scan(&userData).Error
	return userData, err
}
