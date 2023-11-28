package usecase

import (
	"fmt"

	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type OrderUseCase struct {
	orderRepo  interfaces.OrderRepository
	couponRepo interfaces.CouponRepository
}

func NewOrderUseCase(orderRepo interfaces.OrderRepository, couponRepo interfaces.CouponRepository) services.OrderUseCase {
	return &OrderUseCase{
		orderRepo:  orderRepo,
		couponRepo: couponRepo,
	}
}

// OrderAll implements interfaces.OrderUseCase.
func (o *OrderUseCase) OrderAll(id int, paymentTypeId int, CouponName string) (response.ResponseOrder, error) {
	coupon, _ := o.couponRepo.CouponFromName(CouponName)
	if coupon.Id == 0 && CouponName != "" {
		return response.ResponseOrder{}, fmt.Errorf("invalid coupon code")
	}
	order, err := o.orderRepo.OrderAll(id, paymentTypeId, coupon)
	return order, err
}

// UserCancelOrder implements interfaces.OrderUseCase.
func (o *OrderUseCase) UserCancelOrder(orderId int, userId int) error {
	err := o.orderRepo.UserCancelOrder(orderId, userId)
	return err
}

// Displayorder implements interfaces.OrderUseCase.
func (o *OrderUseCase) Displayorder(userId int, orderId int) (response.ResponseOrder, error) {
	order, err := o.orderRepo.DisplayOrder(userId, orderId)
	return order, err
}

// ListAllOrders implements interfaces.OrderUseCase.
func (o *OrderUseCase) ListAllOrders(userId int, queryParams helperStruct.QueryParams) ([]response.OrderResponse, error) {
	orders, err := o.orderRepo.ListAllOrders(userId, queryParams)
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
func (o *OrderUseCase) ListAllOrdersForAdmin(queryParams helperStruct.QueryParams) ([]response.AdminOrder, error) {
	orders, err := o.orderRepo.ListAllOrdersForAdmin(queryParams)
	return orders, err
}

// DisplayOrderForAdmin implements interfaces.OrderUseCase.
func (o *OrderUseCase) DisplayOrderForAdmin(orderId int) (response.AdminOrder, error) {
	order, err := o.orderRepo.DisplayOrderForAdmin(orderId)
	return order, err
}
