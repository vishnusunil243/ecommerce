package interfaces

import "main.go/internal/common/response"

type WishlistUseCase interface {
	AddToWishlist(productId, userId int) error
	RemoveFromWishlist(productId, userId int) error
	ListAllWishlist(userId int) ([]response.Wishlist, error)
	DisplayWishlistProduct(productId, userId int) (response.Wishlist, error)
}
