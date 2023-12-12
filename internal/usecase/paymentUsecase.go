package usecase

import (
	"fmt"

	"github.com/razorpay/razorpay-go"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/domain"
	"main.go/internal/infrastructure/config"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type PaymentUseCase struct {
	paymentRepo interfaces.PaymentRepository
	orderRepo   interfaces.OrderRepository
	cfg         config.Config
}

func NewPaymentuseCase(paymentRepo interfaces.PaymentRepository, orderRepo interfaces.OrderRepository, cfg config.Config) services.PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
		cfg:         cfg,
	}
}

func (c *PaymentUseCase) CreateRazorpayPayment(orderId int) (response.OrderResponse, string, int, error) {
	paymentDetails, err := c.paymentRepo.ViewPaymentDetails(orderId)
	if err != nil {
		return response.OrderResponse{}, "", 0, err
	}

	if paymentDetails.PaymentStatusId == 3 {
		return response.OrderResponse{}, "", 0, fmt.Errorf("payment already completed")
	}
	userId, err := c.orderRepo.UserIdFromOrder(orderId)
	if err != nil {
		return response.OrderResponse{}, "", 0, err
	}
	fmt.Println("user id ", userId)
	//fetch order details from the db
	order, err := c.orderRepo.DisplayOrder(userId, orderId)
	if err != nil {
		return response.OrderResponse{}, "", userId, err
	}
	fmt.Println("order id problem ", order.OrderResponse.Id)
	fmt.Println(order.OrderResponse.Id)
	if order.OrderResponse.Id == 0 {
		return response.OrderResponse{}, "", userId, fmt.Errorf("no such order found")
	}
	client := razorpay.NewClient(c.cfg.RAZORPAYID, c.cfg.RAZORPAYSECRET)

	data := map[string]interface{}{
		"amount":   order.OrderResponse.OrderTotal * 100,
		"currency": "INR",
		"receipt":  "test_receipt_id",
	}

	body, err := client.Order.Create(data, nil)
	if err != nil {
		return response.OrderResponse{}, "", userId, err
	}

	value := body["id"]
	razorpayID := value.(string)
	return order.OrderResponse, razorpayID, userId, err
}

func (c *PaymentUseCase) UpdatePaymentDetails(paymentVerifier helperStruct.PaymentVerification) error {
	paymentDetails, err := c.paymentRepo.ViewPaymentDetails(paymentVerifier.OrderID)
	if err != nil {
		return err
	}
	if paymentDetails.OrdersId == 0 {
		return fmt.Errorf("no order found")
	}

	if paymentDetails.OrderTotal != int(paymentVerifier.Total) {
		return fmt.Errorf("payment amount and order amount does not match")
	}
	updatedPayment, err := c.paymentRepo.UpdatePaymentDetails(paymentVerifier.OrderID, paymentVerifier.PaymentRef)
	if err != nil {
		return err
	}
	if updatedPayment.OrdersId == 0 {
		return fmt.Errorf("failed to update payment details")
	}
	return nil
}

// AddPaymentStatus implements interfaces.PaymentUseCase.
func (p *PaymentUseCase) AddPaymentStatus(paymentStatus helperStruct.PaymentStatus) (domain.PaymentStatus, error) {
	newPaymentStatus, err := p.paymentRepo.AddPaymentStatus(paymentStatus)
	return newPaymentStatus, err
}

// AddPaymentType implements interfaces.PaymentUseCase.
func (p *PaymentUseCase) AddPaymentType(paymentType helperStruct.PaymentType) (domain.PaymentType, error) {
	newPaymentType, err := p.paymentRepo.AddPaymentType(paymentType)
	return newPaymentType, err
}

// ListAllPaymentStatuses implements interfaces.PaymentUseCase.
func (p *PaymentUseCase) ListAllPaymentStatuses() ([]domain.PaymentStatus, error) {
	paymentStatuses, err := p.paymentRepo.ListAllPaymentStatuses()
	return paymentStatuses, err
}

// ListAllPaymentTypes implements interfaces.PaymentUseCase.
func (p *PaymentUseCase) ListAllPaymentTypes() ([]domain.PaymentType, error) {
	paymentTypes, err := p.paymentRepo.ListAllPaymentTypes()
	return paymentTypes, err
}

// UpdatePayemntStatus implements interfaces.PaymentUseCase.
func (p *PaymentUseCase) UpdatePayemntStatus(paymentStatus helperStruct.PaymentStatus, paymentStatusId int) error {
	err := p.paymentRepo.UpdatePaymentStatus(paymentStatus, paymentStatusId)
	return err
}

// UpdatePaymentType implements interfaces.PaymentUseCase.
func (p *PaymentUseCase) UpdatePaymentType(paymentType helperStruct.PaymentType, paymentTypeId int) error {
	err := p.paymentRepo.UpdatePaymentType(paymentType, paymentTypeId)
	return err
}
