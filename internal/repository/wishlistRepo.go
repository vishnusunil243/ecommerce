package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
)

type WishlistDatabase struct {
	DB *gorm.DB
}

func NewWishlistRepo(DB *gorm.DB) interfaces.WishlistRepository {
	return &WishlistDatabase{
		DB: DB,
	}
}

// AddToWishlist implements interfaces.WishlistRepository.
func (w *WishlistDatabase) AddToWishlist(productId int, userId int) error {
	var exists bool
	w.DB.Raw(`SELECT EXISTS (select 1 from wishlists where product_item_id=$1 AND user_id=$2)`, productId, userId).Scan(&exists)
	if exists {
		return fmt.Errorf("this product is already present in the wishlist ")
	}
	addToWishlist := `INSERT INTO wishlists(user_id,product_item_id) VALUES ($1,$2)`
	err := w.DB.Exec(addToWishlist, userId, productId).Error
	return err

}

// RemoveFromWishlist implements interfaces.WishlistRepository.
func (w *WishlistDatabase) RemoveFromWishlist(productId int, userId int) error {
	var exists bool
	w.DB.Raw(`SELECT EXISTS (select 1 from wishlists where product_item_id=$1 AND user_id=$2)`, productId, userId).Scan(&exists)
	if !exists {
		return fmt.Errorf("this product is not present in the wishlist ")
	}
	removeFromWishlist := `DELETE FROM wishlists WHERE product_item_id=$1 AND user_id=$2`
	err := w.DB.Exec(removeFromWishlist, productId, userId).Error
	return err
}

// ListAllWishlist implements interfaces.WishlistRepository.
func (w *WishlistDatabase) ListAllWishlist(userId int) ([]response.Wishlist, error) {
	var wishlists []response.Wishlist
	listAllWishlist := `SELECT wishlists.*, products.product_name,product_items.ram,product_items.storage,products.brand,product_items.color,product_items.graphic_processor,product_items.price AS price_per_unit,product_items.battery
	FROM wishlists LEFT JOIN products on wishlists.product_item_id=products.id
	LEFT JOIN product_items ON wishlists.product_item_id=product_items.id
	WHERE user_id=?`
	err := w.DB.Raw(listAllWishlist, userId).Scan(&wishlists).Error
	return wishlists, err
}

// DisplayWishlistProduct implements interfaces.WishlistRepository.
func (w *WishlistDatabase) DisplayWishlistProduct(productId int, userId int) (response.Wishlist, error) {
	var exists bool
	w.DB.Raw(`SELECT EXISTS (select 1 from wishlists where product_item_id=$1 AND user_id=$2)`, productId, userId).Scan(&exists)
	if !exists {
		return response.Wishlist{}, fmt.Errorf("this product is not present in the wishlist ")
	}
	var wishlist response.Wishlist
	listAllWishlist := `SELECT wishlists.*, products.product_name,product_items.ram,product_items.storage,products.brand,product_items.color,product_items.graphic_processor,product_items.price AS price_per_unit,product_items.battery
	FROM wishlists LEFT JOIN products on wishlists.product_item_id=products.id
	LEFT JOIN product_items ON wishlists.product_item_id=product_items.id
	WHERE user_id=$1 AND wishlists.product_item_id=$2`
	err := w.DB.Raw(listAllWishlist, userId, productId).Scan(&wishlist).Error
	return wishlist, err
}
