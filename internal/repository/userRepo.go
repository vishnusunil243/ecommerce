package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/domain"
	"main.go/internal/repository/interfaces"
)

type userDatabase struct {
	DB *gorm.DB
}

func NewUserRepo(DB *gorm.DB) interfaces.UserRepository {
	return &userDatabase{
		DB: DB,
	}
}

// UserLogin implements interfaces.UserRepository.
func (c *userDatabase) UserLogin(email string) (domain.Users, error) {
	var user domain.Users
	err := c.DB.Raw("SELECT * FROM users WHERE email=?", email).Scan(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

// UserSignUp implements interfaces.UserRepository.
func (c *userDatabase) UserSignUp(user helperStruct.UserReq) (response.UserData, error) {
	var userData response.UserData
	insertQuery := `INSERT INTO users (name,email,mobile,password,created_at)VALUES($1,$2,$3,$4,NOW()) RETURNING id,name,email,mobile`
	err := c.DB.Raw(insertQuery, user.Name, user.Email, user.Mobile, user.Password).Scan(&userData).Error
	return userData, err

}

// AddAdress implements interfaces.UserRepository.
func (c *userDatabase) AddAdress(id int, address helperStruct.Address) (response.Address, error) {
	var newAdress response.Address
	var exists bool
	selectQuery := `SELECT EXISTS (select 1  from addresses WHERE user_id=?)`
	c.DB.Raw(selectQuery, id).Scan(&exists)
	if !exists {
		address.IsDefault = true
	}
	if address.IsDefault { //Change the default address into false
		changeDefault := `UPDATE addresses SET is_default = $1 WHERE users_id=$2 AND is_default=$3`
		err := c.DB.Exec(changeDefault, false, id, true).Error

		if err != nil {
			return newAdress, err
		}
	}
	insertQuery := `INSERT INTO addresses(house_number,users_id,city,district,landmark,pincode,street,is_default) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING house_number,pincode,street,city,district,landmark,is_default `
	err := c.DB.Raw(insertQuery, address.House_number, id, address.City, address.District, address.Landmark, address.Pincode, address.Street, address.IsDefault).Scan(&newAdress).Error
	return newAdress, err
}

// UpdateAddress implements interfaces.UserRepository.
func (c *userDatabase) UpdateAddress(userId, addressId int, address helperStruct.Address) (response.Address, error) {
	var updatedAddress response.Address
	if address.IsDefault { //Change the default address into false
		changeDefault := `UPDATE addresses SET is_default = $1 WHERE users_id=$2 AND is_default=$3`
		err := c.DB.Exec(changeDefault, false, userId, true).Error

		if err != nil {
			return updatedAddress, err
		}
	}
	UpdateQuery := `UPDATE addresses SET house_number=$1,city=$2,district=$3,landmark=$4,pincode=$5,street=$6,is_default=$7 WHERE id=$8 RETURNING house_number,city,district,landmark,pincode,street,is_default`
	err := c.DB.Raw(UpdateQuery, address.House_number, address.City, address.District, address.Landmark, address.Pincode, address.Street, address.IsDefault, addressId).Scan(&updatedAddress).Error
	return updatedAddress, err
}

// DeleteAddress implements interfaces.UserRepository.
func (c *userDatabase) DeleteAddress(addressId int) error {
	var exists bool
	query := `SELECT EXISTS (select 1 exists from addresses where id=?)`
	c.DB.Raw(query, addressId).Scan(&exists)
	if !exists {
		return fmt.Errorf("error deleting address")
	}
	deleteAddress := `DELETE FROM addresses WHERE id=?`
	err := c.DB.Exec(deleteAddress, addressId).Error
	return err
}

// ViewUserProfile implements interfaces.UserRepository.
func (c *userDatabase) ViewUserProfile(id int) (response.UserProfile, error) {
	var userProfile response.UserProfile

	selectProfileQuery := `
		SELECT users.*, addresses.*
		FROM users
		LEFT JOIN addresses ON users.id = addresses.users_id AND addresses.is_default=true
		WHERE users.id = ? 
	`

	err := c.DB.Raw(selectProfileQuery, id).Scan(&userProfile).Error
	if err != nil {
		return userProfile, err
	}
	return userProfile, err
}

// UpdateMobile implements interfaces.UserRepository.
func (c *userDatabase) UpdateMobile(id int, mobile string) (response.UserProfile, error) {
	var userProfile response.UserProfile
	updateQuery := `UPDATE users SET mobile=$1 WHERE id=$2`
	err := c.DB.Exec(updateQuery, mobile, id).Error
	if err != nil {
		return userProfile, err
	}
	selectProfileQuery := `
	SELECT users.*, addresses.*
	FROM users
	LEFT JOIN addresses ON users.id = addresses.users_id
	WHERE users.id = ? AND addresses.is_default=true
`
	err = c.DB.Raw(selectProfileQuery, id).Scan(&userProfile).Error
	return userProfile, err
}

// ChangePassword implements interfaces.UserRepository.
func (c *userDatabase) ChangePassword(id int, password helperStruct.UpdatePassword) (response.UserProfile, error) {
	var userProfile response.UserProfile
	updateQuery := `UPDATE users SET password=$1 WHERE id=$2`
	err := c.DB.Exec(updateQuery, password.NewPassword, id).Error
	if err != nil {
		return userProfile, err
	}
	selectProfileQuery := `
SELECT users.*, addresses.*
FROM users
LEFT JOIN addresses ON users.id = addresses.users_id
WHERE users.id = ? AND addresses.is_default=true
`
	err = c.DB.Raw(selectProfileQuery, id).Scan(&userProfile).Error
	return userProfile, err
}

// RetrieveUserInformation implements interfaces.UserRepository.
func (c *userDatabase) RetrieveUserInformation(id int) (domain.Users, error) {
	var userData domain.Users
	err := c.DB.Raw(`SELECT * FROM users WHERE id=?`, id).Scan(&userData).Error
	return userData, err
}

// ForgotPassword implements interfaces.UserRepository.
func (c *userDatabase) ForgotPassword(newpassword helperStruct.ForgotPassword) error {
	updateQuery := `UPDATE users SET password=$1 WHERE email=$2`
	err := c.DB.Exec(updateQuery, newpassword.NewPassword, newpassword.Email).Error
	return err
}

// ListAllAddresses implements interfaces.UserRepository.
func (c *userDatabase) ListAllAddresses(userId int) ([]response.Address, error) {
	var addresses []response.Address
	selectAddresses := `SELECT * FROM addresses WHERE users_id=?`
	err := c.DB.Raw(selectAddresses, userId).Scan(&addresses).Error
	return addresses, err
}
