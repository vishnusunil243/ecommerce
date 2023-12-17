package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/domain"
)

type PaymentRepository interface {
	ViewPaymentDetails(orderId int) (domain.PaymentDetails, error)
	UpdatePaymentDetails(orderId int, paymentRef string) (domain.PaymentDetails, error)
	AddPaymentType(paymentType helperStruct.PaymentType) (domain.PaymentType, error)
	UpdatePaymentType(paymentType helperStruct.PaymentType, paymentTypeId int) error
	ListAllPaymentTypes() ([]domain.PaymentType, error)
	AddPaymentStatus(paymemtStatus helperStruct.PaymentStatus) (domain.PaymentStatus, error)
	UpdatePaymentStatus(paymentStatus helperStruct.PaymentStatus, paymentStatusId int) error
	ListAllPaymentStatuses() ([]domain.PaymentStatus, error)
}
