package interfaces

import "main.go/internal/common/helperStruct"

type AdminUseCase interface {
	AdminLogin(admin helperStruct.LoginReq) (string, error)
}
