package interfaces

import "main.go/internal/common/response"

type CartRepository interface {
	CreateCart(Id int) error
	AddToCart(productId, userId int) error
	RemoveFromCart(productId, userId int) error
	ListCart(userId int) (response.ViewCart, error)
}
