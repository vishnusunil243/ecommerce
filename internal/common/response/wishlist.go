package response

type Wishlist struct {
	ProductName       string `json:"productname"`
	Brand             string
	Color             string
	Ram               int
	Battery           int
	Storage           int
	Graphic_Processor string
	PricePerUnit      float64
}
