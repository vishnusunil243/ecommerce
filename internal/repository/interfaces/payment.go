package interfaces

import "main.go/internal/domain"

type PaymentRepository interface {
	ViewPaymentDetails(orderId int) (domain.PaymentDetails, error)
	UpdatePaymentDetails(orderId int, paymentRef string) (domain.PaymentDetails, error)
}
