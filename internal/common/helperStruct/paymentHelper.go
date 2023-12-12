package helperStruct

type PaymentVerification struct {
	UserID     int
	OrderID    int
	PaymentRef string
	Total      float64
}
type PaymentType struct {
	Type string
}
type PaymentStatus struct {
	Status string
}
