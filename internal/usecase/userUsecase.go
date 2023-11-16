package usecase

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/infrastructure/config"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type UserUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUsecase(repo interfaces.UserRepository) services.UserUseCase {
	return &UserUseCase{
		userRepo: repo,
	}
}

// UserSignup implements interfaces.UserUseCase.
func (c *UserUseCase) UserSignup(user helperStruct.UserReq) (response.UserData, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.UserData{}, err
	}
	user.Password = string(hash)
	userData, err := c.userRepo.UserSignUp(user)
	return userData, err
}

// UserLogin implements interfaces.UserUseCase.
func (c *UserUseCase) UserLogin(user helperStruct.LoginReq) (string, error) {
	var cfg config.Config
	userData, err := c.userRepo.UserLogin(user.Email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password))
	fmt.Println(err)
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}
	if userData.IsBlocked {
		return "", fmt.Errorf("user is blocked")
	}
	claims := jwt.MapClaims{
		"id":  userData.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.SECRET))
	if err != nil {
		return "", err
	}

	return ss, nil
}

// AddAdress implements interfaces.UserUseCase.
func (cr *UserUseCase) AddAdress(id int, address helperStruct.Address) (response.Address, error) {
	newAdress, err := cr.userRepo.AddAdress(id, address)
	return newAdress, err
}

// UpdateAdress implements interfaces.UserUseCase.
func (cr *UserUseCase) UpdateAddress(addressId int, address helperStruct.Address) (response.Address, error) {
	updatedAdress, err := cr.userRepo.UpdateAddress(addressId, address)
	return updatedAdress, err
}

// ViewUserProfile implements interfaces.UserUseCase.
func (cr *UserUseCase) ViewUserProfile(id int) (response.UserProfile, error) {
	userProfile, err := cr.userRepo.ViewUserProfile(id)
	return userProfile, err
}

// UpdateMobile implements interfaces.UserUseCase.
func (cr *UserUseCase) UpdateMobile(id int, mobile string) (response.UserProfile, error) {
	userProfile, err := cr.userRepo.UpdateMobile(id, mobile)
	return userProfile, err
}

// ChangePassword implements interfaces.UserUseCase.
func (cr *UserUseCase) ChangePassword(id int, password helperStruct.UpdatePassword) (response.UserProfile, error) {
	var userProfile response.UserProfile
	userData, err := cr.userRepo.RetrieveUserInformation(id)
	if err != nil {
		return userProfile, fmt.Errorf("error retrieving user data")
	}
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(password.OldPassword))
	if err != nil {
		return userProfile, fmt.Errorf("the password you have entered is incorrect")
	}
	if password.OldPassword == password.NewPassword {
		return userProfile, fmt.Errorf("old password and new password can't be the same")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return userProfile, fmt.Errorf("error hashing the new password")
	}
	password.NewPassword = string(hash)
	userProfile, err = cr.userRepo.ChangePassword(id, password)
	return userProfile, err
}

// ForgotPassword implements interfaces.UserUseCase.
func (cr *UserUseCase) ForgotPassword(newpassword helperStruct.ForgotPassword) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(newpassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password")
	}
	newpassword.NewPassword = string(hash)
	err = cr.userRepo.ForgotPassword(newpassword)
	return err
}
