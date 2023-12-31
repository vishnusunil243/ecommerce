package helperStruct

import (
	"mime/multipart"
)

type Category struct {
	Name string `json:"name" validate:"required"`
}
type Brand struct {
	Id          uint
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Category_id uint   `json:"category_id" validate:"required"`
}
type Product struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
	Brand       uint   `json:"brand" validate:"required"`
}

type ProductItem struct {
	Product_id        uint    `json:"productid" validate:"required"`
	Sku               string  `json:"sku" validate:"required"`
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
type ImageHelper struct {
	ImageFile     multipart.File
	ImageType     string
	ProductItemId uint
	ImageSize     int64
	ImageData     []byte
}
type SearchProducts struct {
	SearchProducts string
}
