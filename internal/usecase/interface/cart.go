package interfaces

import "main.go/internal/common/response"

type CartUseCase interface {
	CreateCart(Id int) error
	AddToCart(productId, usersId int) error
	RemoveFromCart(productId, userId int) error
	ListCart(userId int) (response.ViewCart, error)
}
