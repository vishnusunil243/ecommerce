package usecase

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type OrderUseCase struct {
	orderRepo interfaces.OrderRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository) services.OrderUseCase {
	return &OrderUseCase{
		orderRepo: orderRepo,
	}
}

// OrderAll implements interfaces.OrderUseCase.
func (o *OrderUseCase) OrderAll(id int, paymentTypeId int) (response.OrderResponse, error) {
	order, err := o.orderRepo.OrderAll(id, paymentTypeId)
	return order, err
}

// UserCancelOrder implements interfaces.OrderUseCase.
func (o *OrderUseCase) UserCancelOrder(orderId int, userId int) error {
	err := o.orderRepo.UserCancelOrder(orderId, userId)
	return err
}

// Displayorder implements interfaces.OrderUseCase.
func (o *OrderUseCase) Displayorder(userId int, orderId int) (response.OrderResponse, error) {
	order, err := o.orderRepo.DisplayOrder(userId, orderId)
	return order, err
}

// ListAllOrders implements interfaces.OrderUseCase.
func (o *OrderUseCase) ListAllOrders(userId int) ([]response.OrderResponse, error) {
	orders, err := o.orderRepo.ListAllOrders(userId)
	return orders, err
}

// ReturnOrder implements interfaces.OrderUseCase.
func (o *OrderUseCase) ReturnOrder(userId int, orderId int) (response.ReturnOrder, error) {
	returnOrder, err := o.orderRepo.ReturnOrder(userId, orderId)
	return returnOrder, err
}

// UpdateOrderStatus implements interfaces.OrderUseCase.
func (o *OrderUseCase) UpdateOrderStatus(updateOrder helperStruct.UpdateOrder) (response.AdminOrder, error) {
	adminOrder, err := o.orderRepo.UpdateOrderStatus(updateOrder)
	return adminOrder, err
}

// ListAllOrdersForAdmin implements interfaces.OrderUseCase.
func (o *OrderUseCase) ListAllOrdersForAdmin() ([]response.AdminOrder, error) {
	orders, err := o.orderRepo.ListAllOrdersForAdmin()
	return orders, err
}
