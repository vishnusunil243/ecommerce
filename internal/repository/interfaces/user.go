package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/domain"
)

type UserRepository interface {
	UserSignUp(user helperStruct.UserReq) (response.UserData, error)
	UserLogin(email string) (domain.Users, error)
}
