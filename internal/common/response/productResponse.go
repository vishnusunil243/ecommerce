package response

type Category struct {
	Id           int
	CategoryName string
}

type Product struct {
	Id           int `json:",omitempty"`
	Name         string
	Description  string
	Brand        string
	CategoryName string
}
type Brand struct {
	Id            int
	Name          string
	Description   string
	Category_id   string
	Category_name string
}

type ProductItem struct {
	Id                uint
	ProductName       string
	Description       string
	Brand             string
	CategoryName      string
	Sku               string
	QtyInStock        int
	Color             string
	Ram               int
	Battery           int
	ScreenSize        float64
	Storage           int
	Graphic_Processor string
	Price             int
	Image             string `json:"image,omitempty"`
}
type ImageResponse struct {
	ID    int    `json:"id"`
	Image []byte `json:"image"`
}
type Image struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
}
