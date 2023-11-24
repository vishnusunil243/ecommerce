package response

import "time"

type AdminData struct {
	Id    int
	Name  string
	Email string
}
type DashBoard struct {
	TotalRevenue      int
	TotalOrders       int
	TotalProductsSold int
	TotalUsers        int
}

type SalesReport struct {
	Name        string
	PaymentType string
	OrderDate   time.Time
	OrderTotal  int
}
type UserReport struct {
	Name        string
	ReportCount uint
}
type ReportInfo struct {
	Username           string
	ReportCount        uint
	ReportedBY         uint
	ReasonForReporting string
}
