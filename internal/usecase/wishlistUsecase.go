package usecase

import (
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type wishlistUseCase struct {
	wishlistRepo interfaces.WishlistRepository
}

func NewWishlistUseCase(wishlistRepo interfaces.WishlistRepository) services.WishlistUseCase {
	return &wishlistUseCase{
		wishlistRepo: wishlistRepo,
	}
}

// AddToWishlist implements interfaces.WishlistUseCase.
func (w *wishlistUseCase) AddToWishlist(productId int, userId int) error {
	err := w.wishlistRepo.AddToWishlist(productId, userId)

	return err
}

// ListAllWishlist implements interfaces.WishlistUseCase.
func (w *wishlistUseCase) ListAllWishlist(userId int) ([]response.Wishlist, error) {
	wishlists, err := w.wishlistRepo.ListAllWishlist(userId)
	return wishlists, err
}

// RemoveFromWishlist implements interfaces.WishlistUseCase.
func (w *wishlistUseCase) RemoveFromWishlist(productId int, userId int) error {
	err := w.wishlistRepo.RemoveFromWishlist(productId, userId)
	return err
}

// DisplayWishlistProduct implements interfaces.WishlistUseCase.
func (w *wishlistUseCase) DisplayWishlistProduct(productId int, userId int) (response.Wishlist, error) {
	wishlist, err := w.wishlistRepo.DisplayWishlistProduct(productId, userId)
	return wishlist, err
}
