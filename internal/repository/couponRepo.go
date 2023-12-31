package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
)

type CouponDatabase struct {
	DB *gorm.DB
}

func NewCouponRepo(DB *gorm.DB) interfaces.CouponRepository {
	return &CouponDatabase{
		DB: DB,
	}
}

// AddCoupon implements interfaces.CouponRepository.
func (c *CouponDatabase) AddCoupon(coupon helperStruct.Coupon) (response.Coupon, error) {
	var newCoupon response.Coupon
	var exists bool
	c.DB.Raw(`SELECT EXISTS (select  1 from coupons where name=?)`, coupon.Name).Scan(&exists)
	if exists {
		return response.Coupon{}, fmt.Errorf("coupon is already present please add a new unique coupon")
	}
	addCoupon := `INSERT INTO coupons(name,quantity,amount,created_at) VALUES ($1,$2,$3,NOW()) RETURNING *`
	err := c.DB.Raw(addCoupon, coupon.Name, coupon.Quantity, coupon.Amount).Scan(&newCoupon).Error
	return newCoupon, err
}

// UpdateCoupon implements interfaces.CouponRepository.
func (c *CouponDatabase) UpdateCoupon(coupon helperStruct.UpdateCoupon) (response.Coupon, error) {
	var exists bool
	c.DB.Raw(`SELECT EXISTS (SELECT 1 FROM coupons WHERE id=?)`, coupon.Id).Scan(&exists)
	if !exists {
		return response.Coupon{}, fmt.Errorf("no coupon found with given id")
	}
	c.DB.Raw(`SELECT EXISTS (select  1 from coupons where name=$1 and id not in ($2))`, coupon.Name, coupon.Id).Scan(&exists)
	if exists {
		return response.Coupon{}, fmt.Errorf("coupon is already present please add a new unique coupon name")
	}
	var updatedCoupon response.Coupon
	updateCoupon := `UPDATE coupons SET name=$1,amount=$2,quantity=$3 WHERE id=$4 RETURNING *`
	err := c.DB.Raw(updateCoupon, coupon.Name, coupon.Amount, coupon.Quantity, coupon.Id).Scan(&updatedCoupon).Error
	return updatedCoupon, err
}

// DisableCoupon implements interfaces.CouponRepository.
func (c *CouponDatabase) DisableCoupon(couponId int) error {
	var exists bool
	c.DB.Raw(`SELECT EXISTS (SELECT 1 FROM coupons WHERE id=?)`, couponId).Scan(&exists)
	if !exists {
		return fmt.Errorf("no coupon found with given id")
	}
	disableCoupon := `UPDATE coupons SET is_disabled=true WHERE id=?`
	err := c.DB.Exec(disableCoupon, couponId).Error
	return err

}

// ListAllCoupons implements interfaces.CouponRepository.
func (c *CouponDatabase) ListAllCoupons(queryParams helperStruct.QueryParams) ([]response.Coupon, int, error) {
	var coupons []response.Coupon
	getAllCoupons := `SELECT * FROM coupons`
	if queryParams.Query != "" && queryParams.Filter != "" {
		getAllCoupons = fmt.Sprintf("%s WHERE LOWER(%s) LIKE '%%%s%%'", getAllCoupons, queryParams.Filter, strings.ToLower(queryParams.Query))
	}
	var count int
	getTotalCount := fmt.Sprintf("SELECT COUNT(*) FROM (%s)", getAllCoupons)
	err := c.DB.Raw(getTotalCount).Scan(&count).Error
	if err != nil {
		return []response.Coupon{}, 0, err
	}
	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			getAllCoupons = fmt.Sprintf("%s ORDER BY %s DESC", getAllCoupons, queryParams.SortBy)
		} else {
			getAllCoupons = fmt.Sprintf("%s ORDER BY %s ASC", getAllCoupons, queryParams.SortBy)
		}
	}
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		getAllCoupons = fmt.Sprintf("%s LIMIT %d OFFSET %d", getAllCoupons, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		getAllCoupons = fmt.Sprintf("%s LIMIT 10 OFFSET 0", getAllCoupons)

	}
	err = c.DB.Raw(getAllCoupons).Scan(&coupons).Error
	return coupons, count, err
}

// DisplayCoupon implements interfaces.CouponRepository.
func (c *CouponDatabase) DisplayCoupon(couponId int) (response.Coupon, error) {
	var exists bool
	c.DB.Raw(`SELECT EXISTS (SELECT 1 FROM coupons WHERE id=?)`, couponId).Scan(&exists)
	if !exists {
		return response.Coupon{}, fmt.Errorf("no coupon found with given id")
	}
	var coupon response.Coupon
	getAllCoupons := `SELECT * FROM coupons WHERE id=?`
	err := c.DB.Raw(getAllCoupons, couponId).Scan(&coupon).Error
	return coupon, err
}
func (c *CouponDatabase) EnableCoupon(couponId int) error {
	var exists bool
	c.DB.Raw(`SELECT EXISTS (SELECT 1 FROM coupons WHERE id=?)`, couponId).Scan(&exists)
	if !exists {
		return fmt.Errorf("no coupon found with given id")
	}
	enableCoupon := `UPDATE coupons SET is_disabled=false WHERE id=?`
	err := c.DB.Exec(enableCoupon, couponId).Error
	return err
}

// CouponFromName implements interfaces.CouponRepository.
func (c *CouponDatabase) CouponFromName(couponName string) (response.Coupon, error) {
	var coupon response.Coupon
	err := c.DB.Raw(`SELECT * FROM coupons WHERE name=?`, couponName).Scan(&coupon).Error
	return coupon, err
}
