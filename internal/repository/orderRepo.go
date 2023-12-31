package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/domain"
	"main.go/internal/repository/interfaces"
)

type orderDatabase struct {
	DB *gorm.DB
}

func NewOrderRepo(DB *gorm.DB) interfaces.OrderRepository {
	return &orderDatabase{
		DB: DB,
	}
}

// OrderAll implements interfaces.OrderRepository.
func (c *orderDatabase) OrderAll(id int, paymentTypeid int, coupon response.Coupon) (response.ResponseOrder, error) {
	tx := c.DB.Begin()
	var cart domain.Carts
	findCart := `SELECT * FROM carts WHERE user_id=?`
	err := tx.Raw(findCart, id).Scan(&cart).Error
	if err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, err
	}
	if cart.Total == 0 {
		setTotal := `UPDATE carts SET total=sub_total WHERE user_id=?`
		err = tx.Exec(setTotal, id).Error
		if err != nil {
			tx.Rollback()
			return response.ResponseOrder{}, err
		}
	}
	if cart.SubTotal == 0 {
		tx.Rollback()
		return response.ResponseOrder{}, fmt.Errorf("there are no items in cart")
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
	discounts ON discounts.brand_id=brands.id AND expiry_date>NOW()
    WHERE 
    ci.carts_id = $1`
	err = tx.Raw(getProductDetails, cart.Id).Scan(&productDetails).Error
	if err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, err
	}
	var totalWithDiscount float64
	var totalDiscountAmount float64
	for _, item := range productDetails {
		if item.DiscountedPrice != 0 {
			item.Total = item.DiscountedPrice
			totalDiscountAmount += item.DiscountPrice

		}
		totalWithDiscount += item.Total
	}
	cart.Total = int(totalWithDiscount)
	var addressId int
	findAddress := `SELECT id FROM addresses WHERE users_id=? AND is_default=true`
	err = tx.Raw(findAddress, id).Scan(&addressId).Error
	if err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, fmt.Errorf("error finding address")
	}
	if addressId == 0 {
		tx.Rollback()
		return response.ResponseOrder{}, fmt.Errorf("please add an address to complete your order")
	}
	if coupon.Id != 0 {
		var exists bool
		tx.Raw(`SELECT EXISTS (select 1 from user_coupons WHERE user_id=$1 AND coupon_id=$2)`, id, coupon.Id).Scan(&exists)
		if coupon.Quantity == 0 {
			tx.Rollback()
			return response.ResponseOrder{}, fmt.Errorf("this coupon has expired")
		}
		if coupon.IsDisabled {
			tx.Rollback()
			return response.ResponseOrder{}, fmt.Errorf("this coupon is disabled ")
		}
		if exists {
			tx.Rollback()
			return response.ResponseOrder{}, fmt.Errorf("can't add this coupon since you have already exhausted it's usage")
		} else {
			addCoupon := `INSERT INTO user_coupons(user_id,coupon_id) VALUES($1,$2)`
			err := tx.Exec(addCoupon, id, coupon.Id).Error
			if err != nil {
				tx.Rollback()
				return response.ResponseOrder{}, fmt.Errorf("can't add this coupon")
			}
			cart.Total = cart.Total - coupon.Amount
			updateCoupons := `UPDATE coupons SET quantity=quantity-1 WHERE id=?`
			err = tx.Exec(updateCoupons, coupon.Id).Error
			if err != nil {
				tx.Rollback()
				return response.ResponseOrder{}, fmt.Errorf("can't add this coupon")
			}
		}
	}
	var order domain.Orders
	insertOrder := `INSERT INTO orders (user_id,order_date,payment_type_id,shipping_address,order_total,order_status_id,payment_status_id) 
	              VALUES ($1,NOW(),$2,$3,$4,$5,$6) RETURNING *`
	err = tx.Raw(insertOrder, id, paymentTypeid, addressId, cart.Total, 1, 1).Scan(&order).Error
	if err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, fmt.Errorf("error placing order")
	}
	var cartItems []helperStruct.CartItems
	cartDetail := `select ci.product_item_id,ci.quantity,pi.price,pi.qty_in_stock,
	            ((pi.price*discount_percent)/100)*(ci.quantity) AS discount_price,
            	(pi.price * ci.quantity)-((pi.price*discount_percent)/100)*(ci.quantity) AS discounted_price
	             from cart_items ci 
	             join product_items pi on ci.product_item_id = pi.id
				 left join products on products.id=pi.id
				 left join brands on brands.brandname=products.brand 
				 left join discounts on discounts.brand_id=brands.id AND expiry_date>NOW()
				 where ci.carts_id=$1`
	err = tx.Raw(cartDetail, cart.Id).Scan(&cartItems).Error
	if err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, err
	}

	//Add the items in the cart into the orderitems one by one
	for _, items := range cartItems {
		//check whether the item is available
		if items.Quantity > items.QtyInStock {
			return response.ResponseOrder{}, fmt.Errorf("out of stock")
		}
		if items.DiscountedPrice != 0 {
			items.Price = int(items.DiscountedPrice)
		}
		insetOrderItems := `INSERT INTO order_items (orders_id,product_item_id,quantity,price) VALUES($1,$2,$3,$4)`
		err = tx.Exec(insetOrderItems, order.Id, items.ProductItemId, items.Quantity, items.Price).Error

		if err != nil {
			tx.Rollback()
			return response.ResponseOrder{}, err
		}
	}
	//Update the cart total
	updateCart := `UPDATE carts SET total=0,sub_total=0,coupon_id=0 WHERE user_id=?`
	err = tx.Exec(updateCart, id).Error
	if err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, err
	}

	//Remove the items from the cart_items
	for _, items := range cartItems {
		removeCartItems := `DELETE FROM cart_items WHERE carts_id =$1 AND product_item_id=$2`
		err = tx.Exec(removeCartItems, cart.Id, items.ProductItemId).Error
		if err != nil {
			tx.Rollback()
			return response.ResponseOrder{}, err
		}
	}

	//Reduce the product qty in stock details
	for _, items := range cartItems {
		updateQty := `UPDATE product_items SET qty_in_stock=product_items.qty_in_stock-$1 WHERE id=$2`
		err = tx.Exec(updateQty, items.Quantity, items.ProductItemId).Error
		if err != nil {
			tx.Rollback()
			return response.ResponseOrder{}, err
		}
	}

	//update the PaymentDetails table with OrdersID, OrderTotal, PaymentTypeID, PaymentStatusID
	createPaymentDetails := `INSERT INTO payment_details
		   (orders_id,
		   order_total,
		   payment_type_id,
		   payment_status_id,
		   updated_at)
		   VALUES($1,$2,$3,$4,NOW())`
	if err = tx.Exec(createPaymentDetails, order.Id, order.OrderTotal, paymentTypeid, 1).Error; err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, err
	}
	if paymentTypeid == 3 {
		var walletAmount int
		getWalletAmount := `SELECT amount FROM wallets WHERE user_id=$1`
		err = tx.Raw(getWalletAmount, id).Scan(&walletAmount).Error
		if err != nil {
			tx.Rollback()
			return response.ResponseOrder{}, fmt.Errorf("error retrieving amount from wallet")
		}
		if walletAmount >= cart.Total {
			insertWalletHistory := `INSERT INTO wallet_histories (recent_transaction,user_id,balance,time) VALUES ($1,$2,$3,NOW())`
			walletHistory := fmt.Sprintf("%d - %d", walletAmount, cart.Total)
			err = tx.Exec(insertWalletHistory, walletHistory, id, (walletAmount - cart.Total)).Error
			if err != nil {
				tx.Rollback()
				return response.ResponseOrder{}, fmt.Errorf("error inserting wallet history")
			}
			updateWallet := `UPDATE wallets SET amount=amount-$1 WHERE user_id=$2`
			err = tx.Exec(updateWallet, cart.Total, id).Error
			if err != nil {
				tx.Rollback()
				return response.ResponseOrder{}, fmt.Errorf("error updating wallet amount")
			}
			updatePaymentStatus := `UPDATE orders SET payment_status_id=5 WHERE id=$1`
			err = tx.Exec(updatePaymentStatus, order.Id).Error
			if err != nil {
				tx.Rollback()
				return response.ResponseOrder{}, fmt.Errorf("error updating payment status")
			}
			updatePaymentDetails := `UPDATE payment_details SET payment_status_id=5 WHERE orders_id=$1`
			err = tx.Exec(updatePaymentDetails, order.Id).Error
			if err != nil {
				tx.Rollback()
				return response.ResponseOrder{}, fmt.Errorf("error updating payment details")
			}
		} else {
			tx.Rollback()
			return response.ResponseOrder{}, fmt.Errorf("you don't have enough amount in the wallet to complete this transaction please choose a different payment method")
		}

	}
	var orderResponse response.OrderResponse
	err = tx.Raw(`SELECT p.type AS payment_type,o.status AS order_status,addresses.*,orders.*,payment_statuses.status AS payment_status,
	order_items.product_item_id,products.product_name
	FROM orders JOIN payment_types p ON  
	p.id=orders.payment_type_id  JOIN order_statuses o ON orders.order_status_id=o.id
	LEFT JOIN payment_statuses ON payment_statuses.id=orders.payment_status_id
	LEFT JOIN order_items ON order_items.orders_id=orders.id
	LEFT JOIN products ON products.id=order_items.product_item_id
	LEFT JOIN addresses ON orders.shipping_address=addresses.id AND is_default=true  WHERE user_id=$1 AND orders.id=$2`, order.UserId, order.Id).Scan(&orderResponse).Error
	if err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, fmt.Errorf("error retrieving order information")
	}
	if coupon.Name != "" {
		orderResponse.CouponCode = coupon.Name
		orderResponse.CouponAmount = -coupon.Amount
		orderResponse.SubTotal = cart.SubTotal
		err := tx.Exec(`UPDATE user_coupons SET order_id=$1 WHERE coupon_id=$2`, order.Id, coupon.Id).Error
		if err != nil {
			tx.Rollback()
			// fmt.Errorf("error updating coupons table")
			return response.ResponseOrder{}, err
		}
	}
	orderResponse.DiscountPrice = int(-totalDiscountAmount)
	var responseOrder response.ResponseOrder
	var orderProducts []response.OrderProduct
	err = tx.Raw(`SELECT order_items.product_item_id,products.product_name,order_items.quantity FROM orders JOIN order_items ON orders.id=order_items.orders_id
	                JOIN products ON order_items.product_item_id=products.id
	                WHERE user_id=$1 AND orders.id=$2`, order.UserId, order.Id).Scan(&orderProducts).Error
	if err != nil {
		return response.ResponseOrder{}, err
	}
	responseOrder.OrderProducts = orderProducts
	responseOrder.OrderResponse = orderResponse
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, err
	}

	return responseOrder, nil
}

// UserCanceOrder implements interfaces.OrderRepository.
func (o *orderDatabase) UserCancelOrder(orderId int, userId int) error {
	tx := o.DB.Begin()
	var items []helperStruct.CartItems
	findProducts := `SELECT product_item_id,quantity FROM order_items WHERE orders_id=?`
	err := tx.Raw(findProducts, orderId).Scan(&items).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting products from cart_items")
	}
	if len(items) == 0 {
		return fmt.Errorf("no order found with the given id")
	}
	for _, item := range items {
		updateProductItem := `UPDATE product_items SET qty_in_stock=qty_in_stock+$1 WHERE id=$2`
		err = tx.Exec(updateProductItem, item.Quantity, item.ProductItemId).Error
		if err != nil {
			tx.Rollback()
			return err
		}

	}
	removeItems := `DELETE FROM order_items WHERE orders_id=$1`
	err = tx.Exec(removeItems, orderId).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var paymentStatusId int
	getPaymentStatus := `SELECT payment_status_id FROM orders WHERE id=$1`
	err = tx.Raw(getPaymentStatus, orderId).Scan(&paymentStatusId).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error retrieving payment status")
	}
	if paymentStatusId == 5 {
		var price int
		getAmount := `SELECT order_total FROM orders WHERE id=$1`
		err = tx.Raw(getAmount, orderId).Scan(&price).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error retrieving total from orders")
		}
		var walletAmount int
		getWalletAmount := `SELECT amount FROM wallets WHERE user_id=$1`
		err = tx.Raw(getWalletAmount, userId).Scan(&walletAmount).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error retrieving amount from wallet")
		}
		insertWalletHistory := `INSERT INTO wallet_histories (recent_transaction,user_id,balance,time) VALUES ($1,$2,$3,NOW())`
		walletHistory := fmt.Sprintf("%d + %d", walletAmount, price)
		err = tx.Exec(insertWalletHistory, walletHistory, userId, (walletAmount + price)).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting wallet history")
		}
		updateWallet := `UPDATE wallets SET amount=amount+$1 WHERE user_id=$2`
		err = tx.Exec(updateWallet, price, userId).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating wallet")
		}
		cancelOrder := `UPDATE orders SET order_status_id=5,payment_status_id=4 WHERE id=$1 AND user_id=$2`
		err = tx.Exec(cancelOrder, orderId, userId).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		updatePaymentDetails := `UPDATE payment_details SET payment_status_id=4 WHERE orders_id=$1`
		err = tx.Exec(updatePaymentDetails, orderId).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating payment details")
		}

	} else {
		cancelOrder := `UPDATE orders SET order_status_id=5,payment_status_id=3 WHERE id=$1 AND user_id=$2`
		err = tx.Exec(cancelOrder, orderId, userId).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		updatePaymentDetails := `UPDATE payment_details SET payment_status_id=3 WHERE orders_id=$1`
		err = tx.Exec(updatePaymentDetails, orderId).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating payment details")
		}
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil

}

// ListAllOrders implements interfaces.OrderRepository.
func (o *orderDatabase) ListAllOrders(userId int, queryParams helperStruct.QueryParams) ([]response.OrderResponse, int, error) {
	var orders []response.OrderResponse
	findOrders := `SELECT p.type AS payment_type,o.status AS order_status,addresses.*,orders.*,payment_statuses.status AS payment_status
	FROM orders JOIN payment_types p ON  
	p.id=orders.payment_type_id LEFT JOIN order_statuses o ON orders.order_status_id=o.id 
	LEFT JOIN addresses ON orders.shipping_address=addresses.id
	LEFT JOIN payment_statuses ON orders.payment_status_id=payment_statuses.id
	WHERE user_id=?`
	if queryParams.Query != "" && queryParams.Filter != "" {
		findOrders = fmt.Sprintf("%s WHERE LOWER(%s) LIKE '%%%s%%'", findOrders, queryParams.Filter, strings.ToLower(queryParams.Query))
	}
	var count int
	getTotalCount := fmt.Sprintf("SELECT COUNT(*) FROM (%s%d)", findOrders[:len(findOrders)-1], userId)
	err := o.DB.Raw(getTotalCount).Scan(&count).Error
	if err != nil {
		return []response.OrderResponse{}, 0, err
	}
	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			findOrders = fmt.Sprintf("%s ORDER BY %s DESC", findOrders, queryParams.SortBy)
		} else {
			findOrders = fmt.Sprintf("%s ORDER BY %s ASC", findOrders, queryParams.SortBy)
		}
	} else {
		findOrders = fmt.Sprintf("%s ORDER BY orders.order_date DESC", findOrders)
	}
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		findOrders = fmt.Sprintf("%s LIMIT %d OFFSET %d", findOrders, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		findOrders = fmt.Sprintf("%s LIMIT 10 OFFSET 0", findOrders)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		findOrders = fmt.Sprintf("%s LIMIT 10 OFFSET 0", findOrders)
	}
	err = o.DB.Raw(findOrders, userId).
		Scan(&orders).Error
	return orders, count, err
}

// DisplayOrder implements interfaces.OrderRepository.
func (o *orderDatabase) DisplayOrder(userId int, orderId int) (response.ResponseOrder, error) {
	var exists bool
	o.DB.Raw(`SELECT EXISTS (select 1 exists from orders where id=$1 and user_id=$2)`, orderId, userId).Scan(&exists)
	if !exists {
		return response.ResponseOrder{}, fmt.Errorf("no such order")
	}
	var order response.OrderResponse
	var orderProducts []response.OrderProduct
	var res response.ResponseOrder
	err := o.DB.Raw(`SELECT p.type AS payment_type,o.status AS order_status,addresses.*,orders.*,payment_statuses.status AS payment_status,order_items.product_item_id AS product_item_id
	,products.product_name,coupons.name AS coupon_code,coupons.amount AS coupon_amount
	FROM orders JOIN payment_types p ON  
	p.id=orders.payment_type_id LEFT JOIN order_statuses o ON orders.order_status_id=o.id
	LEFT JOIN order_items ON orders.id=order_items.orders_id
	LEFT JOIN addresses ON orders.shipping_address=addresses.id  
	LEFT JOIN products ON order_items.product_item_id=products.id 
	LEFT JOIN payment_statuses ON orders.payment_status_id=payment_statuses.id
	LEFT JOIN user_coupons ON orders.id=user_coupons.order_id
	LEFT JOIN coupons ON user_coupons.coupon_id=coupons.id
	WHERE orders.user_id=$1 AND orders.id=$2`, userId, orderId).Scan(&order).Error
	if err != nil {
		return response.ResponseOrder{}, err
	}
	err = o.DB.Raw(`SELECT order_items.product_item_id,products.product_name,order_items.quantity,product_items.price,
	                (discounts.discount_percent/100)*product_items.price AS discount_price
	                FROM orders JOIN order_items ON orders.id=order_items.orders_id
	                JOIN products ON order_items.product_item_id=products.id
					LEFT JOIN brands ON brands.brandname=products.brand
					LEFT JOIN discounts ON discounts.brand_id=brands.id AND expiry_date>NOW()
					LEFT JOIN product_items ON product_items.id=order_items.product_item_id
	                WHERE user_id=$1 AND orders.id=$2`, userId, orderId).Scan(&orderProducts).Error
	var totalPriceWithoutDiscount int
	for _, item := range orderProducts {
		totalPriceWithoutDiscount += (item.Price * item.Quantity)
	}
	order.CouponAmount = -order.CouponAmount
	order.SubTotal = totalPriceWithoutDiscount
	order.DiscountPrice = -(order.SubTotal - order.OrderTotal + order.CouponAmount)
	order.Id = uint(orderId)
	res.OrderProducts = orderProducts
	res.OrderResponse = order
	return res, err
}

// ReturnOrder implements interfaces.OrderRepository.
func (o *orderDatabase) ReturnOrder(userId int, orderId int) (response.ReturnOrder, error) {
	var order domain.Orders
	tx := o.DB.Begin()
	findOrder := `SELECT * FROM orders WHERE id=$1 AND user_id=$2`
	err := tx.Raw(findOrder, orderId, userId).Scan(&order).Error
	if err != nil {
		tx.Rollback()
		return response.ReturnOrder{}, err
	}
	if order.OrderStatusID == 0 {
		tx.Rollback()
		return response.ReturnOrder{}, fmt.Errorf("no such order to return")
	} else if order.OrderStatusID != 4 {
		tx.Rollback()
		return response.ReturnOrder{}, fmt.Errorf("order is not yet delivered")
	}
	updateOrderStatus := `UPDATE orders SET order_status_id=6,payment_status_id=4 WHERE id=$1 AND user_id=$2`
	err = tx.Exec(updateOrderStatus, orderId, userId).Error
	if err != nil {
		tx.Rollback()
		return response.ReturnOrder{}, fmt.Errorf("error updating order status")
	}
	updatePaymentDetails := `UPDATE payment_details SET payment_status_id=4 WHERE orders_id=$1`
	err = tx.Exec(updatePaymentDetails, orderId).Error
	if err != nil {
		return response.ReturnOrder{}, fmt.Errorf("error updating payment_details")
	}
	var walletAmount int
	getWalletAmount := `SELECT amount FROM wallets WHERE user_id=$1`
	err = tx.Raw(getWalletAmount, userId).Scan(&walletAmount).Error
	if err != nil {
		tx.Rollback()
		return response.ReturnOrder{}, fmt.Errorf("error retrieving amount from wallet")
	}
	insertWalletHistory := `INSERT INTO wallet_histories (recent_transaction,user_id,balance,time) VALUES ($1,$2,$3,NOW())`
	walletHistory := fmt.Sprintf("%d + %d", walletAmount, order.OrderTotal)
	err = tx.Exec(insertWalletHistory, walletHistory, userId, (walletAmount + order.OrderTotal)).Error
	if err != nil {
		tx.Rollback()
		return response.ReturnOrder{}, fmt.Errorf("error inserting wallet history")
	}
	updateWallet := `UPDATE wallets SET amount=amount+$1 WHERE user_id=$2`
	err = tx.Exec(updateWallet, order.OrderTotal, userId).Error
	if err != nil {
		tx.Rollback()
		return response.ReturnOrder{}, fmt.Errorf("error refunding the amount")
	}
	var orderResponse response.ReturnOrder
	returnOrder := `SELECT orders.*,order_statuses.status AS order_status FROM orders JOIN order_statuses ON order_statuses.id=orders.order_status_id WHERE orders.id=?`
	err = tx.Raw(returnOrder, orderId).Scan(&orderResponse).Error
	if err != nil {
		tx.Rollback()
		return response.ReturnOrder{}, fmt.Errorf("error retrieving order_info of returned order")
	}
	orderResponse.RefundStatus = "refund successfull"
	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return response.ReturnOrder{}, err
	}
	return orderResponse, err

}

// UpdateOrderStatus implements interfaces.OrderRepository.
func (o *orderDatabase) UpdateOrderStatus(updateOrder helperStruct.UpdateOrder) (response.AdminOrder, error) {
	var exists bool
	o.DB.Raw(`SELECT EXISTS (select 1 exists from orders where id=?)`, updateOrder.OrderId).Scan(&exists)
	if !exists {
		return response.AdminOrder{}, fmt.Errorf("no such order to update")
	}
	if updateOrder.OrderStatusID != 4 {
		updateOrderStatus := `UPDATE orders SET order_status_id=$1 WHERE id=$2 `
		err := o.DB.Exec(updateOrderStatus, updateOrder.OrderStatusID, updateOrder.OrderId).Error
		if err != nil {
			return response.AdminOrder{}, fmt.Errorf("error updating order status")
		}
	} else {
		updateOrderStatus := `UPDATE orders SET order_status_id=$1,payment_status_id=$2 WHERE id=$3`
		err := o.DB.Exec(updateOrderStatus, updateOrder.OrderStatusID, 5, updateOrder.OrderId).Error
		if err != nil {
			return response.AdminOrder{}, fmt.Errorf("error updating order status")
		}
		err = o.DB.Exec(`UPDATE payment_details SET payment_status_id=5 WHERE orders_id=$1`, updateOrder.OrderId).Error
		if err != nil {
			return response.AdminOrder{}, fmt.Errorf("error updating payment_details")
		}
	}
	var adminOrder response.AdminOrder
	selectOrder := `SELECT orders.id AS order_id,orders.payment_type_id AS payment_type_id,order_statuses.status AS order_status,payment_types.type AS payment_type,payment_statuses.status AS payment_status 
	 FROM orders JOIN order_statuses ON orders.order_status_id=order_statuses.id 
	 JOIN payment_types ON orders.payment_type_id=payment_types.id 
	 JOIN payment_statuses ON orders.payment_status_id=payment_statuses.id
	 WHERE orders.id=$1`
	err := o.DB.Raw(selectOrder, updateOrder.OrderId).Scan(&adminOrder).Error
	return adminOrder, err
}

// ListAllOrdersForAdmin implements interfaces.OrderRepository.
func (o *orderDatabase) ListAllOrdersForAdmin(queryParams helperStruct.QueryParams) ([]response.AdminOrder, int, error) {
	var orders []response.AdminOrder
	findOrders := `SELECT orders.id AS order_id,orders.payment_type_id,order_statuses.status AS order_status,payment_types.type AS payment_type,payment_statuses.status AS payment_status
	FROM orders JOIN order_statuses ON orders.order_status_id=order_statuses.id
	JOIN payment_types ON orders.payment_type_id=payment_types.id 
	JOIN payment_statuses ON orders.payment_status_id=payment_statuses.id
     `
	if queryParams.Query != "" && queryParams.Filter != "" {
		findOrders = fmt.Sprintf("%s WHERE LOWER(%s) LIKE '%%%s%%'", findOrders, queryParams.Filter, strings.ToLower(queryParams.Query))
	}
	var count int
	getTotalCount := fmt.Sprintf("SELECT COUNT(*) FROM (%s)", findOrders)
	err := o.DB.Raw(getTotalCount).Scan(&count).Error
	if err != nil {
		return []response.AdminOrder{}, 0, err
	}
	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			findOrders = fmt.Sprintf("%s ORDER BY %s DESC", findOrders, queryParams.SortBy)
		} else {
			findOrders = fmt.Sprintf("%s ORDER BY %s ASC", findOrders, queryParams.SortBy)
		}
	} else {
		findOrders = fmt.Sprintf("%s ORDER BY orders.order_date DESC", findOrders)
	}
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		findOrders = fmt.Sprintf("%s LIMIT %d OFFSET %d", findOrders, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		findOrders = fmt.Sprintf("%s LIMIT 10 OFFSET 0", findOrders)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		findOrders = fmt.Sprintf("%s LIMIT 10 OFFSET 0", findOrders)
	}
	err = o.DB.Raw(findOrders).Scan(&orders).Error
	return orders, count, err
}

// DisplayOrderForAdmin implements interfaces.OrderRepository.
func (o *orderDatabase) DisplayOrderForAdmin(orderId int) (response.AdminOrder, error) {
	var exists bool
	o.DB.Raw(`SELECT EXISTS (select 1 exists from orders where id=$1)`, orderId).Scan(&exists)
	if !exists {
		return response.AdminOrder{}, fmt.Errorf("no such order")
	}
	var order response.AdminOrder
	err := o.DB.Raw(`SELECT orders.id AS order_id,orders.payment_type_id,p.type AS payment_type,o.status AS order_status,addresses.*,payment_statuses.status AS payment_status FROM orders JOIN payment_types p ON  
	p.id=orders.payment_type_id  JOIN order_statuses o ON orders.order_status_id=o.id
	JOIN addresses ON orders.shipping_address=addresses.id AND is_default=true  
	JOIN payment_statuses ON orders.payment_status_id=payment_statuses.id
	WHERE  orders.id=$1`, orderId).Scan(&order).Error
	return order, err
}

// UserIdFromOrder implements interfaces.OrderRepository.
func (o *orderDatabase) UserIdFromOrder(orderId int) (int, error) {
	var userId int
	err := o.DB.Raw(`SELECT user_id FROM orders WHERE id=?`, orderId).Scan(&userId).Error
	return userId, err
}

// AddOrderStatus implements interfaces.OrderRepository.
func (o *orderDatabase) AddOrderStatus(orderStatus helperStruct.OrderStatus) (response.OrderStatus, error) {
	var exists bool
	chechStatusPresent := `SELECT EXISTS (select 1 from order_statuses where status=?)`
	o.DB.Raw(chechStatusPresent, orderStatus.Status).Scan(&exists)
	if exists {
		return response.OrderStatus{}, fmt.Errorf("this status is already present please add a new one")
	}
	var maxId int
	err := o.DB.Raw(`SELECT COALESCE (MAX(id),0) FROM order_statuses`).Scan(&maxId).Error
	if err != nil {
		return response.OrderStatus{}, fmt.Errorf("error retrieving largest id")
	}
	var newOrderStatus response.OrderStatus
	addOrderStatus := `INSERT INTO order_statuses (id,status) VALUES ($1,$2) RETURNING *`
	err = o.DB.Raw(addOrderStatus, maxId+1, orderStatus.Status).Scan(&newOrderStatus).Error
	return newOrderStatus, err
}

// ListAllOrderStatuses implements interfaces.OrderRepository.
func (o *orderDatabase) ListAllOrderStatuses() ([]response.OrderStatus, error) {
	var orderStatuses []response.OrderStatus
	listAllOrderStatuses := `SELECT * FROM order_statuses`
	err := o.DB.Raw(listAllOrderStatuses).Scan(&orderStatuses).Error
	return orderStatuses, err

}

// UpdateOrderStatuses implements interfaces.OrderRepository.
func (o *orderDatabase) UpdateOrderStatuses(orderStatus helperStruct.OrderStatus) (response.OrderStatus, error) {
	var exists bool
	chechStatusPresent := `SELECT EXISTS (select 1 from order_statuses where id=?)`
	o.DB.Raw(chechStatusPresent, orderStatus.Id).Scan(&exists)
	if !exists {
		return response.OrderStatus{}, fmt.Errorf("there is no status with the given id")
	}
	chechStatusPresentWithName := `SELECT EXISTS (select 1 from order_statuses where status=$1 and id not in ($2))`
	o.DB.Raw(chechStatusPresentWithName, orderStatus.Status, orderStatus.Id).Scan(&exists)
	if exists {
		return response.OrderStatus{}, fmt.Errorf("this status is already present please add a new one")
	}
	var updatedOrderStatus response.OrderStatus
	updateOrderStatus := `UPDATE order_statuses SET status=$1 WHERE id=$2 RETURNING *`
	err := o.DB.Raw(updateOrderStatus, orderStatus.Status, orderStatus.Id).Scan(&updatedOrderStatus).Error
	return updatedOrderStatus, err
}
