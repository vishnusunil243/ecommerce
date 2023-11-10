package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type AdminUseCase interface {
	AdminLogin(admin helperStruct.LoginReq) (string, error)
	ListAllUsers() ([]response.UserDetails, error)
	DisplayUser(id int) (response.UserDetails, error)
}
