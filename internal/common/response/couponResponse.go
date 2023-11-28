package response

type Coupon struct {
	Id         uint
	Name       string
	Quantity   int
	Amount     int
	IsDisabled bool
}
