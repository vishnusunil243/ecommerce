package usecase

import (
	"io"

	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type ProductUsecase struct {
	productRepo interfaces.ProductRepository
}

func NewProductUsecase(productRepo interfaces.ProductRepository) services.ProductUsecase {
	return &ProductUsecase{
		productRepo: productRepo,
	}
}

// CreateCategory implements interfaces.ProductUsecase.
func (cr *ProductUsecase) CreateCategory(category helperStruct.Category) (response.Category, error) {
	newCategory, err := cr.productRepo.CreateCategory(category)
	return newCategory, err
}

// UpdateCategory implements interfaces.ProductUsecase.
func (cr *ProductUsecase) UpdateCategory(category helperStruct.Category, id int) (response.Category, error) {
	updatedCategory, err := cr.productRepo.UpdateCategory(category, id)
	return updatedCategory, err
}

// DeleteCategory implements interfaces.ProductUsecase.
func (cr *ProductUsecase) DeleteCategory(id int) error {
	err := cr.productRepo.DeleteCategory(id)
	return err
}

// ListAllCategories implements interfaces.ProductUsecase.
func (cr *ProductUsecase) ListAllCategories() ([]response.Category, error) {
	categories, err := cr.productRepo.ListAllCategories()
	return categories, err
}

// DisplayCategory implements interfaces.ProductUsecase.
func (cr *ProductUsecase) DisplayCategory(id int) (response.Category, error) {
	category, err := cr.productRepo.DisplayCategory(id)
	return category, err
}

// CreateBrand implements interfaces.ProductUsecase.
func (cr *ProductUsecase) CreateBrand(brand helperStruct.Brand) (response.Brand, error) {
	newBrand, err := cr.productRepo.CreateBrand(brand)
	return newBrand, err
}

// UpdatedBrand implements interfaces.ProductUsecase.
func (cr *ProductUsecase) UpdatedBrand(brand helperStruct.Brand, id int) (response.Brand, error) {
	updatedBrand, err := cr.productRepo.UpdateBrand(brand, id)
	return updatedBrand, err
}

// DeleteBrand implements interfaces.ProductUsecase.
func (cr *ProductUsecase) DeleteBrand(id int) error {
	err := cr.productRepo.DeleteBrand(id)
	return err
}

// ListAllBrands implements interfaces.ProductUsecase.
func (cr *ProductUsecase) ListAllBrands(queryParams helperStruct.QueryParams) ([]response.Brand, error) {
	allBrands, err := cr.productRepo.ListAllBrands(queryParams)
	return allBrands, err
}

// DisplayBrand implements interfaces.ProductUsecase.
func (cr *ProductUsecase) DisplayBrand(id int) (response.Brand, error) {
	brand, err := cr.productRepo.DisplayBrand(id)
	return brand, err
}

// AddProduct implements interfaces.ProductUsecase.
func (cr *ProductUsecase) AddProduct(product helperStruct.Product) (response.Product, error) {
	newProduct, err := cr.productRepo.AddProduct(product)
	return newProduct, err
}

// UpdateProduct implements interfaces.ProductUsecase.
func (cr *ProductUsecase) UpdateProduct(product helperStruct.Product, id int) (response.Product, error) {
	updatedProduct, err := cr.productRepo.UpdateProduct(product, id)
	return updatedProduct, err
}

// DeletProduct implements interfaces.ProductUsecase.
func (cr *ProductUsecase) DeletProduct(id int) error {
	err := cr.productRepo.DeleteProduct(id)
	return err
}

// ListAllProducts implements interfaces.ProductUsecase.
func (cr *ProductUsecase) ListAllProducts(queryParams helperStruct.QueryParams) ([]response.Product, error) {
	products, err := cr.productRepo.ListAllProducts(queryParams)
	return products, err
}

// DisplayProduct implements interfaces.ProductUsecase.
func (cr *ProductUsecase) DisplayProduct(id int) (response.Product, error) {
	product, err := cr.productRepo.DisplayProduct(id)
	return product, err
}

// AddProductItem implements interfaces.ProductUsecase.
func (cr *ProductUsecase) AddProductItem(productItem helperStruct.ProductItem) (response.ProductItem, error) {
	newProductItem, err := cr.productRepo.AddProductItem(productItem)
	return newProductItem, err
}

// UpdateProductItem implements interfaces.ProductUsecase.
func (cr *ProductUsecase) UpdateProductItem(id int, productItem helperStruct.ProductItem) (response.ProductItem, error) {
	updatedProductItem, err := cr.productRepo.UpdateProductItem(id, productItem)
	return updatedProductItem, err
}

// ListAllProductItems implements interfaces.ProductUsecase.
func (cr *ProductUsecase) ListAllProductItems(queryParams helperStruct.QueryParams) ([]response.ProductItem, error) {
	productItems, err := cr.productRepo.ListAllProductItems(queryParams)
	return productItems, err
}

// DeleteProductItem implements interfaces.ProductUsecase.
func (cr *ProductUsecase) DeleteProductItem(id int) error {
	err := cr.productRepo.DeleteProductItem(id)
	return err
}

// DisplayProductItem implements interfaces.ProductUsecase.
func (cr *ProductUsecase) DisplayProductItem(id int) (response.ProductItem, error) {
	productItem, err := cr.productRepo.DisplayProductItem(id)
	return productItem, err
}

// ImageUpload implements interfaces.ProductUsecase.
func (cr *ProductUsecase) ImageUpload(image helperStruct.ImageHelper) (response.ImageResponse, error) {

	// Read the file content
	fileBytes, err := io.ReadAll(image.ImageFile)
	if err != nil {
		return response.ImageResponse{}, err
	}
	image.ImageData = fileBytes
	newImage, err := cr.productRepo.UploadImage(image)
	return newImage, err

}

// DeleteImage implements interfaces.ProductUsecase.
func (cr *ProductUsecase) DeleteImage(id int) error {
	err := cr.productRepo.DeleteImage(id)
	return err
}
