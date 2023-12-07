package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
	"main.go/internal/web/handlerUtil"
)

type OrderHandler struct {
	orderUsecase services.OrderUseCase
	adminUsecase services.AdminUseCase
}

func NewOrderHandler(orderUsecase services.OrderUseCase, adminUsecase services.AdminUseCase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
		adminUsecase: adminUsecase,
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
	var CouponName helperStruct.CouponName
	err = c.BindJSON(&CouponName)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	order, err := o.orderUsecase.OrderAll(userId, paymentTypeId, CouponName.CouponName)
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
	var queryParams helperStruct.QueryParams
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.SortBy = c.Query("sort_by")
	queryParams.Filter = c.Query("filter")
	queryParams.Query = c.Query("query")
	if c.Query("sort_desc") != "" {
		queryParams.SortDesc = true
	}
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
	orders, totalCount, err := o.orderUsecase.ListAllOrders(userId, queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving orders",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if len(orders) == 0 {
		c.JSON(http.StatusOK, response.Response{
			StatusCode: 200,
			Message:    "There are no orders yet",
			Data:       nil,
			Errors:     nil,
		})
		return
	}

	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}
	responseStruct := struct {
		Orders    []response.OrderResponse
		NoOfPages int
	}{
		Orders:    orders,
		NoOfPages: totalCount / queryParams.Limit,
	}
	if responseStruct.NoOfPages == 0 {
		responseStruct.NoOfPages = 1
	} else if totalCount%queryParams.Limit != 0 {
		responseStruct.NoOfPages = responseStruct.NoOfPages + 1
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "orders retrieved successfully",
		Data:       responseStruct,
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
	var queryParams helperStruct.QueryParams
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.SortBy = c.Query("sort_by")
	queryParams.Filter = c.Query("filter")
	queryParams.Query = c.Query("query")
	if c.Query("sort_desc") != "" {
		queryParams.SortDesc = true
	}
	orders, totalCount, err := o.orderUsecase.ListAllOrdersForAdmin(queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing all orders",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if len(orders) == 0 {
		c.JSON(http.StatusOK, response.Response{
			StatusCode: 200,
			Message:    "There are no orders yet",
		})
		return
	}
	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}
	responseStruct := struct {
		Orders    []response.AdminOrder
		NoOfPages int
	}{
		Orders:    orders,
		NoOfPages: totalCount / queryParams.Limit,
	}
	if responseStruct.NoOfPages == 0 {
		responseStruct.NoOfPages = 1
	} else if totalCount%queryParams.Limit != 0 {
		responseStruct.NoOfPages = responseStruct.NoOfPages + 1
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "orders listed successfully",
		Data:       responseStruct,
		Errors:     nil,
	})
}
func (o *OrderHandler) DisplayOrderForAdmin(c *gin.Context) {
	paramId := c.Param("order_id")
	orderId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing order id",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	order, err := o.orderUsecase.DisplayOrderForAdmin(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying order",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "order displayed successfully",
		Data:       order,
		Errors:     nil,
	})

}
func (o *OrderHandler) AddOrderStatus(c *gin.Context) {
	var orderStatus helperStruct.OrderStatus
	err := c.BindJSON(&orderStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	newOrderStatus, err := o.orderUsecase.AddOrderStatus(orderStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error adding new orderStatus",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "orderstatus added successfully",
		Data:       newOrderStatus,
		Errors:     nil,
	})
}
func (o *OrderHandler) UpdateOrderStatuses(c *gin.Context) {
	var orderStatus helperStruct.OrderStatus
	err := c.BindJSON(&orderStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	updatedOrderStatus, err := o.orderUsecase.UpdateOrderStatuses(orderStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating order statuses",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "order status updated successfully",
		Data:       updatedOrderStatus,
		Errors:     nil,
	})
}
func (o *OrderHandler) ListAllOrderStatuses(c *gin.Context) {
	orderStatuses, err := o.orderUsecase.ListAllOrderStatuses()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing all orderstatuses",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "orderstatuses listed successfully",
		Data:       orderStatuses,
		Errors:     nil,
	})
}
func (o *OrderHandler) InvoiceDownload(c *gin.Context) {
	paramId := c.Param("order_id")
	orderId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error converting ordeId to integer",
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
	order, err := o.orderUsecase.Displayorder(userId, orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving orderInformation",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if order.OrderResponse.PaymentStatus != "completed" {
		err = fmt.Errorf("please complete the payment to download the invoice")
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Cannot download invoice unless you have completed the payment",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	user, err := o.adminUsecase.DisplayUser(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error loading user details",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	// Generate PDF from JSON data
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(40, 10, "Invoice")
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "pc4u")
	pdf.Ln(8)
	pdf.Cell(40, 10, "Contact : 8129987917")
	pdf.Ln(10)

	// Add order information
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Order Information")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Order ID: %d", order.OrderResponse.Id))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Order Date: %s", order.OrderResponse.OrderDate.Format("2006-01-02 15:04:05")))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Payment Type: %s", order.OrderResponse.PaymentType))
	pdf.Ln(8)

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Buyer Information")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Name: %s ", user.Name))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Email: %s", user.Email))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Phone: %s", user.Mobile))
	pdf.Ln(8)

	// Add shipping address
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Shipping Address")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("House Number: %s", order.OrderResponse.House_number))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Street: %s", order.OrderResponse.Street))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("City: %s", order.OrderResponse.City))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("District: %s", order.OrderResponse.District))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Landmark: %s", order.OrderResponse.Landmark))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Pincode: %d", order.OrderResponse.Pincode))
	pdf.Ln(10)

	// Add order status and payment details
	if order.OrderResponse.CouponCode != "" {
		pdf.Ln(10)
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(40, 10, "Coupon Details")
		pdf.Ln(8)

		pdf.SetFont("Arial", "", 12)
		pdf.Cell(40, 10, fmt.Sprintf("Coupon: %s", order.OrderResponse.CouponCode))
		pdf.Ln(8)
	}

	// Add product details
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Product Details")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	for _, product := range order.OrderProducts {
		pdf.Cell(40, 10, fmt.Sprintf("Product: %s", product.ProductName))
		pdf.Ln(8)
		pdf.Cell(40, 10, fmt.Sprintf("Price: Rs.%d", product.Price))
		if product.DiscountPrice != 0 {
			pdf.Ln(8)
			pdf.Cell(40, 10, fmt.Sprintf("Discount Price: Rs.%.2f", product.DiscountPrice))
		}
		pdf.Ln(8)
		pdf.Cell(40, 10, fmt.Sprintf("Quantity: %d", product.Quantity))
		pdf.Ln(10)

	}

	// Add order summary

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Order Summary")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Subtotal: Rs.%d", order.OrderResponse.SubTotal))
	if order.OrderResponse.CouponCode != "" {
		pdf.Ln(8)
		pdf.Cell(40, 10, fmt.Sprintf("Coupon Amount: Rs.%d", order.OrderResponse.CouponAmount))
	}
	if order.OrderResponse.DiscountPrice != 0 {
		pdf.Ln(8)
		pdf.Cell(40, 10, fmt.Sprintf("Discount Price: Rs.%d", order.OrderResponse.DiscountPrice))
	}
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Order Total: Rs.%d", order.OrderResponse.OrderTotal))
	pdf.Ln(10)
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 10, "Thank you for shopping with us")
	pdf.Ln(10)

	// Set headers for file download
	c.Header("Content-Disposition", "attachment; filename=order.pdf")
	c.Header("Content-Type", "application/pdf")

	// Output the PDF to the response writer
	err = pdf.Output(c.Writer)
	if err != nil {
		// Handle error appropriately
		fmt.Println("Error generating PDF:", err)
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "downloaded response successfully",
		Data:       nil,
		Errors:     nil,
	})
}
