package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type OrderRepository interface {
	OrderAll(UserId, PaymentTypeid int, coupon response.Coupon) (response.ResponseOrder, error)
	UserCancelOrder(orderId, userId int) error
	ListAllOrders(userId int, queryParams helperStruct.QueryParams) ([]response.OrderResponse, int, error)
	DisplayOrder(userId, orderId int) (response.ResponseOrder, error)
	ReturnOrder(userId, orderId int) (response.ReturnOrder, error)
	UpdateOrderStatus(updateOrder helperStruct.UpdateOrder) (response.AdminOrder, error)
	ListAllOrdersForAdmin(queryParams helperStruct.QueryParams) ([]response.AdminOrder, int, error)
	DisplayOrderForAdmin(orderId int) (response.AdminOrder, error)
	UserIdFromOrder(orderId int) (int, error)
	AddOrderStatus(orderStatus helperStruct.OrderStatus) (response.OrderStatus, error)
	UpdateOrderStatuses(orderStatus helperStruct.OrderStatus) (response.OrderStatus, error)
	ListAllOrderStatuses() ([]response.OrderStatus, error)
}
