package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type UserUseCase interface {
	UserSignup(user helperStruct.UserReq) (response.UserData, error)
	UserLogin(user helperStruct.LoginReq) (string, error)
	AddAdress(id int, address helperStruct.Address) (response.Address, error)
	UpdateAddress(addressId int, address helperStruct.Address) (response.Address, error)
	ViewUserProfile(id int) (response.UserProfile, error)
	UpdateMobile(id int, mobile string) (response.UserProfile, error)
	ChangePassword(id int, password helperStruct.UpdatePassword) (response.UserProfile, error)
	ForgotPassword(newpassword helperStruct.ForgotPassword) error
}
