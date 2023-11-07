package repository

import (
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
