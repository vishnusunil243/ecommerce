package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
)

type CouponHandler struct {
	couponUsecase services.CouponUsecase
}

func NewCouponHandler(couponUsecase services.CouponUsecase) *CouponHandler {
	return &CouponHandler{
		couponUsecase: couponUsecase,
	}
}
func (cu *CouponHandler) AddCoupon(c *gin.Context) {
	var coupon helperStruct.Coupon
	err := c.BindJSON(&coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	newCoupon, err := cu.couponUsecase.AddCoupon(coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error adding coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "coupon added successfully",
		Data:       newCoupon,
		Errors:     nil,
	})
}
func (cu *CouponHandler) UpdateCoupon(c *gin.Context) {
	var coupon helperStruct.UpdateCoupon
	err := c.BindJSON(&coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	paramId := c.Param("coupon_id")
	couponId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	coupon.Id = couponId
	updatedCoupon, err := cu.couponUsecase.UpdateCoupon(coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "coupon updated successfully",
		Data:       updatedCoupon,
		Errors:     nil,
	})
}
func (cu *CouponHandler) DisableCoupon(c *gin.Context) {
	paramId := c.Param("coupon_id")
	couponId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = cu.couponUsecase.DisableCoupon(couponId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error disabling coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "coupon disabled successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (cu *CouponHandler) ListAllCoupons(c *gin.Context) {
	coupns, err := cu.couponUsecase.ListAllCoupons()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing all coupons",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "coupons listed successfully",
		Data:       coupns,
		Errors:     nil,
	})
}
func (cu *CouponHandler) DisplayCoupon(c *gin.Context) {
	paramId := c.Param("coupon_id")
	couponId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	coupon, err := cu.couponUsecase.DisplayCoupon(couponId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "coupon displayed successfully",
		Data:       coupon,
		Errors:     nil,
	})
}
func (cu *CouponHandler) EnableCoupon(c *gin.Context) {
	paramId := c.Param("coupon_id")
	couponId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error enabling coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = cu.couponUsecase.EnableCoupon(couponId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error enabling coupon",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "coupon enabled successfully",
		Data:       nil,
		Errors:     nil,
	})
}
