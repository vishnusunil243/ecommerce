package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type PaymentUseCase interface {
	CreateRazorpayPayment(orderId int) (response.OrderResponse, string, int, error)
	UpdatePaymentDetails(paymentVerifier helperStruct.PaymentVerification) error
}
