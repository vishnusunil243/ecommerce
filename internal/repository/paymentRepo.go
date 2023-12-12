package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
	"main.go/internal/domain"
	"main.go/internal/repository/interfaces"
)

type PaymentDatabase struct {
	DB *gorm.DB
}

func NewPaymentRepo(DB *gorm.DB) interfaces.PaymentRepository {
	return &PaymentDatabase{DB}
}

func (c *PaymentDatabase) ViewPaymentDetails(orderID int) (domain.PaymentDetails, error) {
	var paymentDetails domain.PaymentDetails
	fetchPaymentDetailsQuery := `SELECT * FROM payment_details WHERE orders_id = $1;`
	err := c.DB.Raw(fetchPaymentDetailsQuery, orderID).Scan(&paymentDetails).Error
	fmt.Println("2", paymentDetails)
	return paymentDetails, err
}

func (c *PaymentDatabase) UpdatePaymentDetails(orderID int, paymentRef string) (domain.PaymentDetails, error) {
	var updatedPayment domain.PaymentDetails
	var updatePaymentQuery string
	if paymentRef != "" {
		updatePaymentQuery = `	UPDATE payment_details SET payment_type_id = 2, payment_status_id = 5, payment_ref = $1, updated_at = NOW()
							WHERE orders_id = $2 RETURNING *;`
	} else {
		updatePaymentQuery = `UPDATE payment_details SET payment_type_id = 2,payment_status_id = 6,updated_at = NOW()
		                       WHERE orders_id=$2 RETURNING *`
	}
	tx := c.DB.Begin()
	err := tx.Raw(updatePaymentQuery, paymentRef, orderID).Scan(&updatedPayment).Error
	if err != nil {
		tx.Rollback()
		return updatedPayment, err
	}
	var updateOrderTable string
	if paymentRef != "" {
		updateOrderTable = `UPDATE orders SET payment_status_id=5 WHERE id=$1`
	} else {
		updateOrderTable = `UPDATE orders SET payment_status_id=6 WHERE id=$1`
	}
	err = tx.Exec(updateOrderTable, orderID).Error
	if err != nil {
		tx.Rollback()
		return domain.PaymentDetails{}, err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return domain.PaymentDetails{}, err
	}
	return updatedPayment, err
}

// AddPaymentStatus implements interfaces.PaymentRepository.
func (p *PaymentDatabase) AddPaymentStatus(paymemtStatus helperStruct.PaymentStatus) (domain.PaymentStatus, error) {
	var exists bool
	p.DB.Raw(`SELECT EXISTS (select 1 from payment_statuses where status=$1)`, paymemtStatus.Status).Scan(&exists)
	if exists {
		return domain.PaymentStatus{}, fmt.Errorf("this payment status already exists please add a new unique one or update an existing payment status")
	}
	var maxId int
	err := p.DB.Raw(`SELECT COALESCE (MAX(id),0) FROM payment_statuses`).Scan(&maxId).Error
	if err != nil {
		return domain.PaymentStatus{}, err
	}
	var newPaymentStatus domain.PaymentStatus
	addPaymentStatus := `INSERT INTO payment_statuses(id,status) VALUES ($1,$2) RETURNING *`
	err = p.DB.Raw(addPaymentStatus, (maxId + 1), paymemtStatus.Status).Scan(&newPaymentStatus).Error
	return newPaymentStatus, err
}

// AddPaymentType implements interfaces.PaymentRepository.
func (p *PaymentDatabase) AddPaymentType(paymentType helperStruct.PaymentType) (domain.PaymentType, error) {
	var exists bool
	p.DB.Raw(`SELECT EXISTS (select 1 from payment_types where type=$1)`, paymentType.Type).Scan(&exists)
	if exists {
		return domain.PaymentType{}, fmt.Errorf("payment type already exists")
	}
	var maxId int
	err := p.DB.Raw(`SELECT COALESCE(MAX(id),0) FROM payment_types`).Scan(&maxId).Error
	if err != nil {
		return domain.PaymentType{}, err
	}
	addPaymentType := `INSERT INTO payment_types VALUES($1,$2) RETURNING *`
	var newePaymentType domain.PaymentType
	err = p.DB.Raw(addPaymentType, (maxId + 1), paymentType.Type).Scan(&newePaymentType).Error
	return newePaymentType, err
}

// ListAllPaymentStatuses implements interfaces.PaymentRepository.
func (p *PaymentDatabase) ListAllPaymentStatuses() ([]domain.PaymentStatus, error) {
	var paymentStatuses []domain.PaymentStatus
	err := p.DB.Raw(`SELECT * FROM  payment_statuses`).Scan(&paymentStatuses).Error
	return paymentStatuses, err
}

// ListAllPaymentTypes implements interfaces.PaymentRepository.
func (p *PaymentDatabase) ListAllPaymentTypes() ([]domain.PaymentType, error) {
	var paymentTypes []domain.PaymentType
	err := p.DB.Raw(`SELECT * FROM payment_types`).Scan(&paymentTypes).Error
	return paymentTypes, err
}

// UpdatePaymentStatus implements interfaces.PaymentRepository.
func (p *PaymentDatabase) UpdatePaymentStatus(paymentStatus helperStruct.PaymentStatus, paymentStatusId int) error {
	var exists bool
	p.DB.Raw(`SELECT EXISTS (select 1 from payment_statuses where id=?)`, paymentStatusId).Scan(&exists)
	if !exists {
		return fmt.Errorf("no payment status found with the given id")
	}
	updatePaymentStatus := `UPDATE payment_statuses SET status=$1 WHERE id=$2`
	err := p.DB.Exec(updatePaymentStatus, paymentStatus.Status, paymentStatusId).Error
	return err
}

// UpdatePaymentType implements interfaces.PaymentRepository.
func (p *PaymentDatabase) UpdatePaymentType(paymentType helperStruct.PaymentType, paymentTypeId int) error {
	var exists bool
	p.DB.Raw(`SELECT EXISTS (select 1 from payment_types where id=?)`, paymentTypeId).Scan(&exists)
	if !exists {
		return fmt.Errorf("no payment type found with the given id")
	}
	updatePaymentType := `UPDATE payment_types SET type=$1 WHERE id=$2`
	err := p.DB.Exec(updatePaymentType, paymentType.Type, paymentTypeId).Error
	return err
}
