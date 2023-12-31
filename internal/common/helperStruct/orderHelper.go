package helperStruct

type Cart struct {
	Id     int
	Tottal int
}

type CartItems struct {
	ProductItemId   int
	Quantity        int
	Price           int
	DiscountPrice   float64
	DiscountedPrice float64
	QtyInStock      int
}

type UpdateOrder struct {
	OrderId       uint
	OrderStatusID uint
}
type OrderStatus struct {
	Id     uint   `json:"id,omitempty"`
	Status string `json:"status"`
}
