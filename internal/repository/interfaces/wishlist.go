package interfaces

import "main.go/internal/common/response"

type WishlistRepository interface {
	AddToWishlist(productId int, userId int) error
	RemoveFromWishlist(productId int, userId int) error
	ListAllWishlist(userId int) ([]response.Wishlist, error)
	DisplayWishlistProduct(productId int, userId int) (response.Wishlist, error)
}
