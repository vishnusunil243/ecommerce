package interfaces

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
)

type ProductUsecase interface {
	CreateCategory(category helperStruct.Category) (response.Category, error)
	UpdateCategory(category helperStruct.Category, id int) (response.Category, error)
	DeleteCategory(id int) error
	ListAllCategories() ([]response.Category, error)
	DisplayCategory(id int) (response.Category, error)
	CreateBrand(brand helperStruct.Brand) (response.Brand, error)
	UpdatedBrand(brand helperStruct.Brand, id int) (response.Brand, error)
	DeleteBrand(id int) error
	ListAllBrands(queryParams helperStruct.QueryParams) ([]response.Brand, int, error)
	DisplayBrand(id int) (response.Brand, error)
	AddProduct(product helperStruct.Product) (response.Product, error)
	UpdateProduct(product helperStruct.Product, id int) (response.Product, error)
	DeletProduct(id int) error
	ListAllProducts(queryParams helperStruct.QueryParams) ([]response.Product, int, error)
	DisplayProduct(id int) (response.Product, error)
	AddProductItem(productItem helperStruct.ProductItem) (response.ProductItem, error)
	UpdateProductItem(id int, productItem helperStruct.ProductItem) (response.ProductItem, error)
	ListAllProductItems(queryParams helperStruct.QueryParams) ([]response.ProductItem, int, error)
	DeleteProductItem(id int) error
	// ImageUpload(image helperStruct.ImageHelper) (response.ImageResponse, error)
	UploadImage(filepath string, productid int) (response.Image, error)
	DeleteImage(id int) error
	DisplayProductItem(id int) (response.DisplayProductItem, error)
	SearchProducts(queryParams helperStruct.QueryParams, searchProducts string) ([]response.Product, error)
}
