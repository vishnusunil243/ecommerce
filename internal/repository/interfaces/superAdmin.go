package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/domain"
)

type SuperAdminRepository interface {
	Login(superadmin helperStruct.SuperLoginReq) (domain.SuperAdmin, error)
	CreateAdmin(admin helperStruct.CreateAdmin) (response.AdminData, error)
	ListAllAdmins() ([]response.AdminData, error)
	DisplayAdmin(id int) (response.AdminData, error)
	BlockAdmin(id int) (response.AdminData, error)
	BlockUser(id int) (response.UserData, error)
}
