package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type OrderUseCase interface {
	OrderAll(id, paymentTypeId int, CouponName string) (response.ResponseOrder, error)
	UserCancelOrder(orderId, userId int) error
	ListAllOrders(userId int, queryParams helperStruct.QueryParams) ([]response.OrderResponse, error)
	Displayorder(userId, orderId int) (response.ResponseOrder, error)
	ReturnOrder(userId, orderId int) (response.ReturnOrder, error)
	UpdateOrderStatus(updateOrder helperStruct.UpdateOrder) (response.AdminOrder, error)
	ListAllOrdersForAdmin(queryParams helperStruct.QueryParams) ([]response.AdminOrder, error)
	DisplayOrderForAdmin(orderId int) (response.AdminOrder, error)
}
