package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type CouponUsecase interface {
	AddCoupon(coupon helperStruct.Coupon) (response.Coupon, error)
	UpdateCoupon(coupon helperStruct.UpdateCoupon) (response.Coupon, error)
	DisableCoupon(couponId int) error
	ListAllCoupons(queryParams helperStruct.QueryParams) ([]response.Coupon, int, error)
	DisplayCoupon(couponId int) (response.Coupon, error)
	EnableCoupon(couponId int) error
}
