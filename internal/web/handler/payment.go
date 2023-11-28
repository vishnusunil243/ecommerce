package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
)

type PaymentHandler struct {
	paymentUseCase services.PaymentUseCase
}

func NewPaymentHandler(paymentUseCase services.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
	}
}

func (cr *PaymentHandler) CreateRazorpayPayment(c *gin.Context) {
	paramsId := c.Param("orderId")
	orderId, err := strconv.Atoi(paramsId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "cant find order id",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	fmt.Println(paramsId)

	order, razorpayID, userId, err := cr.paymentUseCase.CreateRazorpayPayment(orderId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Can't complete order",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.HTML(200, "app.html", gin.H{
		"UserID":       userId,
		"total_price":  order.OrderTotal,
		"total":        order.OrderTotal,
		"orderData":    order.Id,
		"orderid":      razorpayID,
		"amount":       order.OrderTotal,
		"Email":        "vishnusunil243@gmail.com",
		"Phone_Number": "8129987917",
	})
}

func (cr *PaymentHandler) PaymentSuccess(c *gin.Context) {

	paymentRef := c.Query("payment_ref")
	fmt.Println("paymentRef from query :", paymentRef)

	idStr := c.Query("order_id")
	fmt.Print("order id from query _:", idStr)

	idStr = strings.ReplaceAll(idStr, " ", "")

	orderID, err := strconv.Atoi(idStr)
	fmt.Println("_converted order  id from query :", orderID)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Can't find orderId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	uID := c.Query("user_id")
	userID, err := strconv.Atoi(uID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Can't find UserId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	t := c.Query("total")
	fmt.Println("total from query :", t)
	total, err := strconv.ParseFloat(t, 32)
	fmt.Println("total from query converted:", total)

	if err != nil {
		//	handle err
		fmt.Println("failed to fetch order id")
	}

	//orderID := strings.Trim("orderid", " ")

	paymentVerifier := helperStruct.PaymentVerification{
		UserID:     userID,
		OrderID:    orderID,
		PaymentRef: paymentRef,
		Total:      total,
	}

	fmt.Println(paymentVerifier.OrderID)

	fmt.Println("payment verifier in handler : ", paymentVerifier)
	//paymentVerifier.
	err = cr.paymentUseCase.UpdatePaymentDetails(paymentVerifier)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "faild to update payment",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "payment updated",
		Data:       nil,
		Errors:     nil,
	})
}
