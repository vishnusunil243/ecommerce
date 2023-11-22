package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/domain"
)

type UserRepository interface {
	UserSignUp(user helperStruct.UserReq) (response.UserData, error)
	UserLogin(email string) (domain.Users, error)
	AddAdress(id int, address helperStruct.Address) (response.Address, error)
	UpdateAddress(userId, addressId int, addess helperStruct.Address) (response.Address, error)
	DeleteAddress(addressId, userId int) error
	ListAllAddresses(userId int) ([]response.Address, error)
	ViewUserProfile(id int) (response.UserProfile, error)
	UpdateMobile(id int, mobile string) (response.UserProfile, error)
	RetrieveUserInformation(id int) (domain.Users, error)
	ChangePassword(id int, password helperStruct.UpdatePassword) (response.UserProfile, error)
	ForgotPassword(newpassword helperStruct.ForgotPassword) error
}
