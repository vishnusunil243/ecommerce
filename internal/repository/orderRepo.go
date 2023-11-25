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
func (c *orderDatabase) OrderAll(id int, paymentTypeid int) (response.ResponseOrder, error) {
	tx := c.DB.Begin()
	var cart domain.Carts
	findCart := `SELECT * FROM carts WHERE user_id=?`
	err := tx.Raw(findCart, id).Scan(&cart).Error
	if err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, err
	}
	if cart.Total == 0 {
		setTotal := `UPDATE carts SET total=sub_total WHERE users_id=?`
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
	var order domain.Orders
	insertOrder := `INSERT INTO orders (user_id,order_date,payment_type_id,shipping_address,order_total,order_status_id,payment_status_id) 
	              VALUES ($1,NOW(),$2,$3,$4,$5,$6) RETURNING *`
	err = tx.Raw(insertOrder, id, paymentTypeid, addressId, cart.Total, 1, 1).Scan(&order).Error
	if err != nil {
		tx.Rollback()
		return response.ResponseOrder{}, fmt.Errorf("error placing order")
	}
	var cartItems []helperStruct.CartItems
	cartDetail := `select ci.product_item_id,ci.quantity,pi.price,pi.qty_in_stock  from cart_items ci join product_items pi on ci.product_item_id = pi.id where ci.carts_id=$1`
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
	var orderResponse response.OrderResponse
	err = tx.Raw(`SELECT orders.*,p.type AS payment_type,o.status AS order_status,addresses.*,payment_statuses.status AS payment_status,
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
	cancelOrder := `UPDATE orders SET order_status_id=5,payment_status_id=3 WHERE id=$1 AND user_id=$2`
	err = tx.Exec(cancelOrder, orderId, userId).Error
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

// ListAllOrders implements interfaces.OrderRepository.
func (o *orderDatabase) ListAllOrders(userId int, queryParams helperStruct.QueryParams) ([]response.OrderResponse, error) {
	var orders []response.OrderResponse
	findOrders := `SELECT orders.*,p.type AS payment_type,o.status AS order_status,addresses.*,payment_statuses.status AS payment_status
	,order_items.product_item_id AS product_item_id,products.product_name
	FROM orders JOIN payment_types p ON  
	p.id=orders.payment_type_id LEFT JOIN order_statuses o ON orders.order_status_id=o.id 
	LEFT JOIN addresses ON orders.shipping_address=addresses.id
	LEFT JOIN order_items ON orders.id=order_items.orders_id
	LEFT JOIN products ON order_items.product_item_id=products.id 
	LEFT JOIN payment_statuses ON orders.payment_status_id=payment_statuses.id
	WHERE user_id=?`
	if queryParams.Query != "" && queryParams.Filter != "" {
		findOrders = fmt.Sprintf("%s WHERE LOWER(%s) LIKE '%%%s%%'", findOrders, queryParams.Filter, strings.ToLower(queryParams.Query))
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
	err := o.DB.Raw(findOrders, userId).
		Scan(&orders).Error
	return orders, err
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
	err := o.DB.Raw(`SELECT orders.*,p.type AS payment_type,o.status AS order_status,addresses.*,payment_statuses.status AS payment_status,order_items.product_item_id AS product_item_id
	,products.product_name
	FROM orders JOIN payment_types p ON  
	p.id=orders.payment_type_id LEFT JOIN order_statuses o ON orders.order_status_id=o.id
	LEFT JOIN order_items ON orders.id=order_items.orders_id
	LEFT JOIN addresses ON orders.shipping_address=addresses.id  
	LEFT JOIN products ON order_items.product_item_id=products.id 
	LEFT JOIN payment_statuses ON orders.payment_status_id=payment_statuses.id
	WHERE user_id=$1 AND orders.id=$2`, userId, orderId).Scan(&order).Error
	if err != nil {
		return response.ResponseOrder{}, err
	}
	err = o.DB.Raw(`SELECT order_items.product_item_id,products.product_name,order_items.quantity FROM orders JOIN order_items ON orders.id=order_items.orders_id
	                JOIN products ON order_items.product_item_id=products.id
	                WHERE user_id=$1 AND orders.id=$2`, userId, orderId).Scan(&orderProducts).Error
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
func (o *orderDatabase) ListAllOrdersForAdmin(queryParams helperStruct.QueryParams) ([]response.AdminOrder, error) {
	var orders []response.AdminOrder
	findOrders := `SELECT orders.id AS order_id,orders.payment_type_id,order_statuses.status AS order_status,payment_types.type AS payment_type,payment_statuses.status AS payment_status
	FROM orders JOIN order_statuses ON orders.order_status_id=order_statuses.id
	JOIN payment_types ON orders.payment_type_id=payment_types.id 
	JOIN payment_statuses ON orders.payment_status_id=payment_statuses.id
     `
	if queryParams.Query != "" && queryParams.Filter != "" {
		findOrders = fmt.Sprintf("%s WHERE LOWER(%s) LIKE '%%%s%%'", findOrders, queryParams.Filter, strings.ToLower(queryParams.Query))
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
	err := o.DB.Raw(findOrders).Scan(&orders).Error
	return orders, err
}

// DisplayOrderForAdmin implements interfaces.OrderRepository.
func (o *orderDatabase) DisplayOrderForAdmin(orderId int) (response.AdminOrder, error) {
	var exists bool
	o.DB.Raw(`SELECT EXISTS (select 1 exists from orders where id=$1)`, orderId).Scan(&exists)
	if !exists {
		return response.AdminOrder{}, fmt.Errorf("no such order")
	}
	var order response.AdminOrder
	err := o.DB.Raw(`SELECT orders.*,p.type AS payment_type,o.status AS order_status,addresses.*,payment_statuses.status AS payment_status FROM orders JOIN payment_types p ON  
	p.id=orders.payment_type_id  JOIN order_statuses o ON orders.order_status_id=o.id
	JOIN addresses ON orders.shipping_address=addresses.id AND is_default=true  
	JOIN payment_statuses ON orders.payment_status_id=payment_statuses.id
	WHERE  orders.id=$1`, orderId).Scan(&order).Error
	return order, err
}
