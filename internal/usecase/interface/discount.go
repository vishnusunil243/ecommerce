package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type DiscountUseCase interface {
	AddDiscount(discount helperStruct.Discount) (response.Discount, error)
	UpdateDiscount(discount helperStruct.Discount, discountId uint) (response.Discount, error)
	ListAllDiscount(queryParams helperStruct.QueryParams) ([]response.Discount, int, error)
	DeleteDiscount(id int) error
}
