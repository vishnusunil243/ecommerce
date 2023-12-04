package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type DiscountRepository interface {
	AddDiscount(discount helperStruct.Discount) (response.Discount, error)
	UpdateDiscount(discount helperStruct.Discount, discountId uint) (response.Discount, error)
	DeleteDiscount(id int) error
	ListAllDiscount() ([]response.Discount, error)
}
