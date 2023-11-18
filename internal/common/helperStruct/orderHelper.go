package helperStruct

type Cart struct {
	Id     int
	Tottal int
}

type CartItems struct {
	ProductItemId int
	Quantity      int
	Price         int
	QtyInStock    int
}

type UpdateOrder struct {
	OrderId       uint
	OrderStatusID uint
}
