package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
	"main.go/internal/web/handlerUtil"
)

type OrderHandler struct {
	orderUsecase services.OrderUseCase
}

func NewOrderHandler(orderUsecase services.OrderUseCase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}
func (o *OrderHandler) OrderAll(c *gin.Context) {
	paramId := c.Param("payment_id")
	paymentTypeId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing paymenttype",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving userid",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	order, err := o.orderUsecase.OrderAll(userId, paymentTypeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "can't place order",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "order placed",
		Data:       order,
		Errors:     nil,
	})
}
func (o *OrderHandler) UserCancelOrder(c *gin.Context) {
	paramId := c.Param("order_id")
	orderId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving order id",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving userId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = o.orderUsecase.UserCancelOrder(orderId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error cancelling order",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "order cancelled successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (o *OrderHandler) ListAllOrders(c *gin.Context) {
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving user id from contex",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	orders, err := o.orderUsecase.ListAllOrders(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving orders",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "orders retrieved successfully",
		Data:       orders,
		Errors:     nil,
	})
}
func (o *OrderHandler) DisplayOrder(c *gin.Context) {
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving user id from context",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	paramId := c.Param("order_id")
	orderId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	order, err := o.orderUsecase.Displayorder(userId, orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving order info",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "order retrieved successfully",
		Data:       order,
		Errors:     nil,
	})
}
func (o *OrderHandler) ReturnOrder(c *gin.Context) {
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving userId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	paramId := c.Param("order_id")
	orderId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving orderId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	returnOrder, err := o.orderUsecase.ReturnOrder(userId, orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "order can't be returned",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "order returned successfully",
		Data:       returnOrder,
		Errors:     nil,
	})
}
func (o *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	var updateOrder helperStruct.UpdateOrder
	err := c.BindJSON(&updateOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	updatedOrder, err := o.orderUsecase.UpdateOrderStatus(updateOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating order status",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "order status updated successfully",
		Data:       updatedOrder,
		Errors:     nil,
	})

}
func (o *OrderHandler) ListAllOrdersForAdmin(c *gin.Context) {
	orders, err := o.orderUsecase.ListAllOrdersForAdmin()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing all orders",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "orders listed successfully",
		Data:       orders,
		Errors:     nil,
	})
}
