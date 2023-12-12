package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/domain"
)

type PaymentUseCase interface {
	CreateRazorpayPayment(orderId int) (response.OrderResponse, string, int, error)
	UpdatePaymentDetails(paymentVerifier helperStruct.PaymentVerification) error
	AddPaymentType(paymentType helperStruct.PaymentType) (domain.PaymentType, error)
	UpdatePaymentType(paymentType helperStruct.PaymentType, paymentTypeId int) error
	ListAllPaymentTypes() ([]domain.PaymentType, error)
	AddPaymentStatus(paymentStatus helperStruct.PaymentStatus) (domain.PaymentStatus, error)
	UpdatePayemntStatus(paymentStatus helperStruct.PaymentStatus, paymentStatusId int) error
	ListAllPaymentStatuses() ([]domain.PaymentStatus, error)
}
