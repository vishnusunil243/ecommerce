package usecase

import (
	"fmt"

	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type CouponUsecase struct {
	couponRepo interfaces.CouponRepository
}

func NewCouponUsecase(couponRepo interfaces.CouponRepository) services.CouponUsecase {
	return &CouponUsecase{
		couponRepo: couponRepo,
	}
}

// AddCoupon implements interfaces.CouponUsecase.
func (c *CouponUsecase) AddCoupon(coupon helperStruct.Coupon) (response.Coupon, error) {
	if coupon.Amount < 0 {
		return response.Coupon{}, fmt.Errorf("amount can't have a negative value")
	}
	if coupon.Quantity < 0 {
		return response.Coupon{}, fmt.Errorf("quantity can't have a negative value")
	}
	newCoupon, err := c.couponRepo.AddCoupon(coupon)
	return newCoupon, err
}

// UpdateCoupon implements interfaces.CouponUsecase.
func (c *CouponUsecase) UpdateCoupon(coupon helperStruct.UpdateCoupon) (response.Coupon, error) {
	if coupon.Amount < 0 {
		return response.Coupon{}, fmt.Errorf("amount can't have a negative value")
	}
	if coupon.Quantity < 0 {
		return response.Coupon{}, fmt.Errorf("quantity can't have a negative value")
	}
	updatedCoupon, err := c.couponRepo.UpdateCoupon(coupon)
	return updatedCoupon, err
}

// DisableCoupon implements interfaces.CouponUsecase.
func (c *CouponUsecase) DisableCoupon(couponId int) error {
	err := c.couponRepo.DisableCoupon(couponId)
	return err
}

// ListAllCoupons implements interfaces.CouponUsecase.
func (c *CouponUsecase) ListAllCoupons(queryParams helperStruct.QueryParams) ([]response.Coupon, int, error) {
	coupons, totalCount, err := c.couponRepo.ListAllCoupons(queryParams)
	return coupons, totalCount, err
}

// DisplayCoupon implements interfaces.CouponUsecase.
func (c *CouponUsecase) DisplayCoupon(couponId int) (response.Coupon, error) {
	coupon, err := c.couponRepo.DisplayCoupon(couponId)
	return coupon, err
}

// EnableCoupon implements interfaces.CouponUsecase.
func (c *CouponUsecase) EnableCoupon(couponId int) error {
	err := c.couponRepo.EnableCoupon(couponId)
	return err
}
