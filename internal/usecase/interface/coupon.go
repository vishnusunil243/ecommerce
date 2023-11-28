package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type CouponUsecase interface {
	AddCoupon(coupon helperStruct.Coupon) (response.Coupon, error)
	UpdateCoupon(coupon helperStruct.UpdateCoupon) (response.Coupon, error)
	DisableCoupon(couponId int) error
	ListAllCoupons() ([]response.Coupon, error)
	DisplayCoupon(couponId int) (response.Coupon, error)
	EnableCoupon(couponId int) error
}
