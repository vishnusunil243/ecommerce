package response

type DisplayCart struct {
	ProductName       string `json:"productname"`
	Brand             string
	Color             string
	Ram               int
	Battery           int
	Storage           int
	Graphic_Processor string
	Quantity          int
	PricePerUnit      float64
	Total             float64
}
type ViewCart struct {
	CartItems []DisplayCart `json:"cart_items"`
	SubTotal  float64       `json:"sub_total"`
	CartTotal float64       `json:"cart_total"`
}
