package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/domain"
)

type AdminRepository interface {
	AdminLogin(email string) (domain.Admins, error)
	ListAllUsers(queryParams helperStruct.QueryParams) ([]response.UserDetails, error)
	DispalyUser(id int) (response.UserDetails, error)
	ReportUser(usersid int) (response.UserReport, error)
	GetDashBoard(dashboard helperStruct.Dashboard) (response.DashBoard, error)
}
