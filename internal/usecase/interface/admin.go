package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type AdminUseCase interface {
	AdminLogin(admin helperStruct.LoginReq) (string, error)
	ListAllUsers(queryParams helperStruct.QueryParams) ([]response.UserDetails, int, error)
	DisplayUser(id int) (response.UserDetails, error)
	ReportUser(UsersId int) (response.UserReport, error)
	GetDashboard(dashboard helperStruct.Dashboard) (response.DashBoard, error)
	ViewSalesReport(filter helperStruct.Dashboard) ([]response.SalesReport, error)
}
