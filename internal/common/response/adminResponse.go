package response

import "time"

type AdminData struct {
	Id           int
	Name         string
	Email        string
	IsSuperAdmin bool
}
type DashBoard struct {
	TotalRevenue        int
	TotalOrders         int
	TotalProductsSelled int
	TotalUsers          int
}

type SalesReport struct {
	Name        string
	PaymentType string
	OrderDate   time.Time
	OrderTotal  int
}
