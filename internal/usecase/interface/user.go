package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type UserUseCase interface {
	UserSignup(user helperStruct.UserReq) (response.UserData, error)
	UserLogin(user helperStruct.LoginReq) (string, error)
}
