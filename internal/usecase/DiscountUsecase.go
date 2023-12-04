package usecase

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type DiscouneUseCase struct {
	discountRepo interfaces.DiscountRepository
}

func NewDiscountUseCase(dicountRepo interfaces.DiscountRepository) services.DiscountUseCase {
	return &DiscouneUseCase{
		discountRepo: dicountRepo,
	}
}

// AddDiscount implements interfaces.DiscountUseCase.
func (d *DiscouneUseCase) AddDiscount(discount helperStruct.Discount) (response.Discount, error) {
	newDiscount, err := d.discountRepo.AddDiscount(discount)
	return newDiscount, err
}

// DeleteDiscount implements interfaces.DiscountUseCase.
func (d *DiscouneUseCase) DeleteDiscount(id int) error {
	err := d.discountRepo.DeleteDiscount(id)
	return err
}

// ListAllDiscount implements interfaces.DiscountUseCase.
func (d *DiscouneUseCase) ListAllDiscount() ([]response.Discount, error) {
	discounts, err := d.discountRepo.ListAllDiscount()
	return discounts, err
}

// UpdateDiscount implements interfaces.DiscountUseCase.
func (d *DiscouneUseCase) UpdateDiscount(discount helperStruct.Discount, discountId uint) (response.Discount, error) {
	updatedDiscount, err := d.discountRepo.UpdateDiscount(discount, discountId)
	return updatedDiscount, err
}
