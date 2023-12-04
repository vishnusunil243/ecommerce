package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/response"
	"main.go/internal/domain"
	"main.go/internal/repository/interfaces"
)

type cartDatabase struct {
	DB *gorm.DB
}

func NewCartRepo(DB *gorm.DB) interfaces.CartRepository {
	return &cartDatabase{
		DB: DB,
	}
}

// CreateCart implements interfaces.CartRepository.
func (c *cartDatabase) CreateCart(Id int) error {
	insertQuery := `INSERT INTO carts(user_id,sub_total,total,coupon_id) VALUES ($1,0,0,0)`
	err := c.DB.Exec(insertQuery, Id).Error
	return err
}

// AddToCart implements interfaces.CartRepository.
func (c *cartDatabase) AddToCart(productId int, userId int) error {
	tx := c.DB.Begin()
	var cartId int
	findCartId := `SELECT id FROM carts WHERE user_id=?`
	err := tx.Raw(findCartId, userId).Scan(&cartId).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var cartItemId int
	cartItemCheck := `SELECT id FROM cart_items WHERE carts_id=$1 AND product_item_id=$2 LIMIT 1`
	err = tx.Raw(cartItemCheck, cartId, productId).Scan(&cartItemId).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if cartItemId == 0 {
		addToCart := `INSERT INTO cart_items (carts_id,product_item_id,quantity)VALUES($1,$2,1)`
		err = tx.Exec(addToCart, cartId, productId).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		updatCart := `UPDATE cart_items SET quantity = cart_items.quantity+1 WHERE id = $1 `
		err = tx.Exec(updatCart, cartItemId).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	//finding the price of the product
	var price int
	findPrice := `SELECT price FROM product_items WHERE id=$1`
	err = tx.Raw(findPrice, productId).Scan(&price).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	//Updating the subtotal in cart table
	var subtotal int
	updateSubTotal := `UPDATE carts SET sub_total=carts.sub_total+$1 WHERE user_id=$2 RETURNING sub_total`
	err = tx.Raw(updateSubTotal, price, userId).Scan(&subtotal).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	updateTotal := `UPDATE carts SET total=$1 where id=$2`
	err = tx.Exec(updateTotal, subtotal, cartId).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil

}

// RemoveFromCart implements interfaces.CartRepository.
func (c *cartDatabase) RemoveFromCart(productId int, userId int) error {
	tx := c.DB.Begin()

	//Find cart id
	var cartId int
	findCartId := `SELECT id FROM carts WHERE user_id=? `
	err := tx.Raw(findCartId, userId).Scan(&cartId).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	//Find the qty of the product in cart
	var qty int
	findQty := `SELECT quantity FROM cart_items WHERE carts_id=$1 AND product_item_id=$2`
	err = tx.Raw(findQty, cartId, productId).Scan(&qty).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if qty == 0 {
		tx.Rollback()
		return fmt.Errorf("no items in cart to reomve")
	}

	//If the qty is 1 dlt the product from the cart
	if qty == 1 {
		dltItem := `DELETE FROM cart_items WHERE carts_id=$1 AND product_item_id=$2`
		err := tx.Exec(dltItem, cartId, productId).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else { // If there is  more than one product reduce the qty by 1
		updateQty := `UPDATE cart_items SET quantity=cart_items.quantity-1 WHERE carts_id=$1 AND product_item_id=$2`
		err = tx.Exec(updateQty, cartId, productId).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	//Find the price of the product item
	var price int
	productPrice := `SELECT price FROM product_items WHERE id=$1`
	err = tx.Raw(productPrice, productId).Scan(&price).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	//Update the subtotal reduce the price of the cart total with price of the product
	var subTotal int
	updateSubTotalAndTotal := `UPDATE carts SET sub_total=sub_total-$1,total=total-$2 WHERE user_id=$3 RETURNING sub_total`
	err = tx.Raw(updateSubTotalAndTotal, price, price, userId).Scan(&subTotal).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil

}

// ListCart implements interfaces.CartRepository.
func (c *cartDatabase) ListCart(userId int) (response.ViewCart, error) {
	tx := c.DB.Begin()
	type cartDetails struct {
		Id       int
		SubTotal float64
		Total    float64
	}
	var cart cartDetails
	getCartDetails := `SELECT 
	c.id,
	c.sub_total,
	c.total
	FROM carts c WHERE c.user_id=$1`
	err := tx.Raw(getCartDetails, userId).Scan(&cart).Error
	if err != nil {
		tx.Rollback()
		return response.ViewCart{}, err
	}
	var cartItems domain.CartItem
	getCartItemDetails := `SELECT * FROM cart_items WHERE carts_id=$1`
	err = tx.Raw(getCartItemDetails, cart.Id).Scan(&cartItems).Error
	if err != nil {
		tx.Rollback()
		return response.ViewCart{}, err
	}
	var productDetails []response.DisplayCart
	getProductDetails := `SELECT 
    p.brand,
    pr.product_name,
    pi.sku AS product_sku,
    pi.color,
    pi.ram,
    pi.battery,
    pi.graphic_processor,
    pi.storage,
    ci.quantity,
    pi.price AS price_per_unit,
    (pi.price * ci.quantity) AS total ,
	((pi.price*discount_percent)/100)*(ci.quantity) AS discount_price,
	(pi.price * ci.quantity)-((pi.price*discount_percent)/100)*(ci.quantity) AS discounted_price
   FROM 
    cart_items ci 
    JOIN 
    product_items pi ON ci.product_item_id = pi.id
    JOIN 
    products pr ON pi.product_id = pr.id
    JOIN 
    products p ON pi.product_id = p.id 
	LEFT JOIN
	brands ON brands.brandname=p.brand
	LEFT JOIN 
	discounts ON discounts.brand_id=brands.id
    WHERE 
    ci.carts_id = $1`
	err = tx.Raw(getProductDetails, cart.Id).Scan(&productDetails).Error
	if err != nil {
		tx.Rollback()
		return response.ViewCart{}, err
	}
	var carts response.ViewCart
	carts.CartTotal = cart.Total
	carts.SubTotal = cart.SubTotal
	carts.CartItems = productDetails
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return response.ViewCart{}, err
	}
	return carts, nil
}
