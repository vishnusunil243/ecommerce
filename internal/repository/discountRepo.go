package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
)

type DiscountDatabase struct {
	DB *gorm.DB
}

func NewDiscountRepo(DB *gorm.DB) interfaces.DiscountRepository {
	return &DiscountDatabase{
		DB: DB,
	}
}

// AddDiscount implements interfaces.DiscountRepository.
func (d *DiscountDatabase) AddDiscount(discount helperStruct.Discount) (response.Discount, error) {
	var exists bool
	d.DB.Raw(`SELECT EXISTS (select 1 from discounts where brand_id=?)`, discount.BrandId).Scan(&exists)
	if exists {
		return response.Discount{}, fmt.Errorf("this brand already has a discount please add discount for a different brand")
	}
	var maxId int
	err := d.DB.Raw(`SELECT COALESCE(MAX(id),0) FROM discounts`).Scan(&maxId).Error
	if err != nil {
		return response.Discount{}, fmt.Errorf("error retrieving maxId")
	}
	var newDiscount response.Discount
	addDiscount := `INSERT INTO discounts(id,discount_percent,brand_id,max_discount_amount,min_purchase_amount,expiry_date) VALUES($1,$2,$3,$4,$5,$6)RETURNING *`
	err = d.DB.Exec(addDiscount, maxId+1, discount.DiscountPercent, discount.BrandId, discount.MaxDiscountAmount, discount.MinPurchaseAmount, discount.ExpiryDate).Error
	if err != nil {
		return response.Discount{}, err
	}
	displayDiscount := `SELECT discounts.*,brands.brandname AS brand_name FROM discounts LEFT JOIN brands ON discounts.brand_id=brands.id WHERE discounts.brand_id=?`
	err = d.DB.Raw(displayDiscount, discount.BrandId).Scan(&newDiscount).Error
	return newDiscount, err

}

// DeleteDiscount implements interfaces.DiscountRepository.
func (d *DiscountDatabase) DeleteDiscount(id int) error {
	var exists bool
	d.DB.Raw(`SELECT EXISTS (select 1 from discounts where id=?)`, id).Scan(&exists)
	if !exists {
		return fmt.Errorf("no such discount to delete")
	}
	deleteDiscount := `DELETE FROM discounts WHERE id=?`
	err := d.DB.Exec(deleteDiscount, id).Error
	return err
}

// ListAllDiscount implements interfaces.DiscountRepository.
func (d *DiscountDatabase) ListAllDiscount() ([]response.Discount, error) {
	var discounts []response.Discount
	listAllDiscount := `SELECT discounts.*,brands.brandname AS brand_name FROM discounts LEFT JOIN brands ON discounts.brand_id=brands.id`
	err := d.DB.Raw(listAllDiscount).Scan(&discounts).Error
	return discounts, err
}

// UpdateDiscount implements interfaces.DiscountRepository.
func (d *DiscountDatabase) UpdateDiscount(discount helperStruct.Discount, discountId uint) (response.Discount, error) {
	var exists bool
	d.DB.Raw(`SELECT EXISTS (select 1 from discounts where id=?)`, discountId).Scan(&exists)
	if !exists {
		return response.Discount{}, fmt.Errorf("no discount found with the given id")
	}
	var updatedDiscount response.Discount
	updateDiscount := `UPDATE discounts SET max_discount_amount=$1,min_purchase_amount=$2,expiry_date=$3,discount_percent=$4,brand_id=$5 WHERE id=$6`
	err := d.DB.Exec(updateDiscount, discount.MaxDiscountAmount, discount.MinPurchaseAmount, discount.ExpiryDate, discount.DiscountPercent, discount.BrandId, discountId).Error
	if err != nil {
		return response.Discount{}, err
	}
	displayDiscount := `SELECT discounts.*,brands.brandname AS brand_name FROM discounts LEFT JOIN brands ON discounts.brand_id=brands.id WHERE discounts.id=?`
	err = d.DB.Raw(displayDiscount, discountId).Scan(&updatedDiscount).Error
	return updatedDiscount, err
}
