package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type OrderUseCase interface {
	OrderAll(id, paymentTypeId int) (response.OrderResponse, error)
	UserCancelOrder(orderId, userId int) error
	ListAllOrders(userId int) ([]response.OrderResponse, error)
	Displayorder(userId, orderId int) (response.OrderResponse, error)
	ReturnOrder(userId, orderId int) (response.ReturnOrder, error)
	UpdateOrderStatus(updateOrder helperStruct.UpdateOrder) (response.AdminOrder, error)
	ListAllOrdersForAdmin() ([]response.AdminOrder, error)
}
