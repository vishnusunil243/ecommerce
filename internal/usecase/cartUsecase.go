package usecase

import (
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type cartUseCase struct {
	cartRepo interfaces.CartRepository
}

func NewCartUseCase(cartRepo interfaces.CartRepository) services.CartUseCase {
	return &cartUseCase{
		cartRepo: cartRepo,
	}
}

// CreateCart implements interfaces.CartUseCase.
func (c *cartUseCase) CreateCart(Id int) error {
	err := c.cartRepo.CreateCart(Id)
	return err
}

// AddToCart implements interfaces.CartUseCase.
func (c *cartUseCase) AddToCart(productId int, usersId int) error {
	err := c.cartRepo.AddToCart(productId, usersId)
	return err
}

// RemoveFromCart implements interfaces.CartUseCase.
func (c *cartUseCase) RemoveFromCart(productId int, userId int) error {
	err := c.cartRepo.RemoveFromCart(productId, userId)
	return err
}

// ListCart implements interfaces.CartUseCase.
func (c *cartUseCase) ListCart(userId int) (response.ViewCart, error) {
	viewCart, err := c.cartRepo.ListCart(userId)
	return viewCart, err
}
