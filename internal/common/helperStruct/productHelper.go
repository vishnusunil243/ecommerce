package helperStruct

type Category struct {
	Name string `json:"name" validate:"required"`
}
type Brand struct {
	Id            uint
	Name          string `json:"name"`
	Description   string `json:"description"`
	Category_name string `json:"category_name"`
}
type Product struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Brand       string `json:"brand" validate:"required"`
}

type ProductItem struct {
	Product_id        uint    `json:"productid"`
	Sku               string  `json:"sku"`
	Qty               int     `json:"quantity"`
	Color             string  `json:"colour"`
	Ram               int     `json:"ram"`
	Battery           int     `json:"battery"`
	Screen_size       float64 `json:"screensize"`
	Storage           int     `json:"storage"`
	Graphic_Processor string  `json:"graphic_processor"`
	Price             int     `json:"price"`
	Image             string  `json:"image"`
}
type QueryParams struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Query    string `json:"query"`   //search key word
	Filter   string `json:"filter"`  //to specify the column name
	SortBy   string `json:"sort_by"` //to specify column to set the sorting
	SortDesc bool   `json:"sort_desc"`
}
