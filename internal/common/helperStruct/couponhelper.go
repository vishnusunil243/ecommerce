package helperStruct

type Coupon struct {
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
	Quantity int    `json:"quantity"`
}
type UpdateCoupon struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
	Quantity int    `json:"quantity"`
}
type CouponName struct {
	CouponName string
}
