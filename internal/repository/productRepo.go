package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
)

type ProductDatabase struct {
	DB *gorm.DB
}

// CreateProduct implements interfaces.ProductRepository.

func NewProductRepo(DB *gorm.DB) interfaces.ProductRepository {
	return &ProductDatabase{
		DB: DB,
	}
}

// CreateCategory implements interfaces.ProductRepository.
func (c *ProductDatabase) CreateCategory(category helperStruct.Category) (response.Category, error) {
	var newcategory response.Category
	var exists bool
	selectQuery := `select exists(select 1 from categories where category_name=$1)`
	c.DB.Raw(selectQuery, category.Name).Scan(&exists)
	if exists {
		return newcategory, fmt.Errorf("category already exists")
	}
	query := `INSERT INTO categories(category_name,created_at) VALUES($1,NOW()) RETURNING id,category_name`
	err := c.DB.Raw(query, category.Name).Scan(&newcategory).Error
	if err != nil {
		return newcategory, err
	}
	return newcategory, nil
}

// ProductCategory implements interfaces.ProductRepository.
func (c *ProductDatabase) UpdateCategory(category helperStruct.Category, id int) (response.Category, error) {
	var updatedCategory response.Category
	updateQuery := `UPDATE categories SET category_name=$1,updated_at=NOW() WHERE id=$2 RETURNING id,category_name `
	err := c.DB.Raw(updateQuery, category.Name, id).Scan(&updatedCategory).Error
	if err != nil {
		return updatedCategory, err
	}
	if updatedCategory.Id == 0 {
		return updatedCategory, fmt.Errorf("no such category to update")
	}
	return updatedCategory, nil
}

// DeleteCategory implements interfaces.ProductRepository.
func (c *ProductDatabase) DeleteCategory(id int) error {
	var exists bool
	query := `select exists(select 1 from categories where id=?)`
	err := c.DB.Raw(query, id).Scan(&exists).Error
	if !exists {
		return err
	}

	errs := c.DB.Exec(`DELETE FROM categories WHERE id=?`, id).Error
	return errs
}

// ListAllCategories implements interfaces.ProductRepository.
func (c *ProductDatabase) ListAllCategories() ([]response.Category, error) {
	var categories []response.Category
	err := c.DB.Raw(`SELECT * FROM categories`).Scan(&categories).Error
	return categories, err
}

// DisplayCategory implements interfaces.ProductRepository.
func (c *ProductDatabase) DisplayCategory(id int) (response.Category, error) {
	var category response.Category
	var exists bool
	query := `select exists(select 1 from categories where id=$1)`
	c.DB.Raw(query, id).Scan(&exists)
	if !exists {
		return category, fmt.Errorf("no such category")
	}
	err := c.DB.Raw(`SELECT * FROM categories WHERE id=?`, id).Scan(&category).Error
	return category, err
}

// CreateBrand implements interfaces.ProductRepository.
func (c *ProductDatabase) CreateBrand(brand helperStruct.Brand) (response.Brand, error) {
	var newbrand response.Brand
	var exists bool
	c.DB.Raw(`select exists(select 1 from brands where brandname=?)`, brand.Name).Scan(&exists)
	if exists {
		return newbrand, fmt.Errorf("brand is already present")
	}
	insertQuery := `INSERT INTO brands (brandname,description,category_id,created_at) VALUES ($1,$2,$3,NOW()) RETURNING id,brandname AS name,description,category_id`
	err := c.DB.Raw(insertQuery, brand.Name, brand.Description, brand.Category_id).Scan(&newbrand).Error
	if err != nil {
		return newbrand, err
	}
	selectQuery := `SELECT category_name FROM categories WHERE id=?`
	err = c.DB.Raw(selectQuery, brand.Category_id).Scan(&newbrand.Category_name).Error
	if err != nil {
		return newbrand, fmt.Errorf("error retrieving category_name")
	}
	return newbrand, nil

}

// UpdateBrand implements interfaces.ProductRepository.
func (c *ProductDatabase) UpdateBrand(brand helperStruct.Brand, id int) (response.Brand, error) {
	var updatedBrand response.Brand

	updateQuery := `UPDATE brands SET brandname=$1,description=$2,category_id=$3,updated_at=NOW() WHERE id=$4 RETURNING id,brandname AS name,category_id,description`
	err := c.DB.Raw(updateQuery, brand.Name, brand.Description, brand.Category_id, id).Scan(&updatedBrand).Error
	if err != nil {
		return updatedBrand, fmt.Errorf("error updating brand")
	}
	selectQuery := `SELECT category_name FROM categories WHERE id=?`
	err = c.DB.Raw(selectQuery, brand.Category_id).Scan(&updatedBrand.Category_name).Error
	if err != nil {
		return updatedBrand, fmt.Errorf("error retrieving category_name")
	}
	if updatedBrand.Id == 0 {
		return updatedBrand, fmt.Errorf("no such brand to update")
	}
	return updatedBrand, err
}

// DeleteBrands implements interfaces.ProductRepository.
func (c *ProductDatabase) DeleteBrand(id int) error {
	var exists bool
	c.DB.Raw(`select exists(select 1 from brands where id=?)`, id).Scan(&exists)
	if !exists {
		return fmt.Errorf("no brand found with given id")
	}
	err := c.DB.Exec(`DELETE FROM brands WHERE id=?`, id).Error
	return err
}

// ListAllBrands implements interfaces.ProductRepository.
func (c *ProductDatabase) ListAllBrands(queryParams helperStruct.QueryParams) ([]response.Brand, error) {
	var brands []response.Brand
	getBrands := `
    SELECT brands.brandname AS name,brands.id,brands.category_id,brands.description, categories.category_name
    FROM brands
    JOIN categories ON brands.category_id = categories.id
	
`
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		getBrands = fmt.Sprintf("%s LIMIT %d OFFSET %d", getBrands, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		getBrands = fmt.Sprintf("%s LIMIT 10 OFFSET 0", getBrands)
	}
	err := c.DB.Raw(getBrands).Scan(&brands).Error
	return brands, err
}

// DisplayBrand implements interfaces.ProductRepository.
func (c *ProductDatabase) DisplayBrand(id int) (response.Brand, error) {
	var brand response.Brand
	var exists bool
	c.DB.Raw(`select exists(select 1 from categories where id=?)`, id).Scan(&exists)
	if !exists {
		return brand, fmt.Errorf("no brand found with given id")
	}

	err := c.DB.Raw(` SELECT brands.brandname AS name,brands.id,brands.category_id,brands.description, categories.category_name
    FROM brands
    JOIN categories ON brands.category_id = categories.id WHERE brands.id=?`, id).Scan(&brand).Error
	return brand, err
}
func (c *ProductDatabase) AddProduct(product helperStruct.Product) (response.Product, error) {
	var brand response.Brand
	var newProduct response.Product
	var exists bool
	c.DB.Raw(`select exists(select 1 from products where product_name=?)`, product.Name).Scan(&exists)
	if exists {
		return newProduct, fmt.Errorf(" product is already present")
	}
	err := c.DB.Raw(`
    SELECT b.brandname AS name,b.category_id, c.category_name
    FROM brands b
    JOIN categories c ON b.category_id = c.id
    WHERE b.id = ?
`, product.Brand).Scan(&brand).Error
	fmt.Println(brand.Name)
	if err != nil {
		return newProduct, err
	}
	insertQuery := `INSERT INTO products (product_name,description,brand,category_id,created_at) VALUES ($1,$2,$3,$4,NOW())
	RETURNING id,product_name AS name,description,brand,category_id`
	err = c.DB.Raw(insertQuery, product.Name, product.Description, brand.Name, brand.Category_id).Scan(&newProduct).Error
	if err != nil {
		return newProduct, err
	}
	newProduct.CategoryName = brand.Category_name
	return newProduct, err
}

// UpdateProducts implements interfaces.ProductRepository.
func (c *ProductDatabase) UpdateProduct(product helperStruct.Product, id int) (response.Product, error) {
	var updatedProduct response.Product
	var exists bool
	c.DB.Raw(`select exists(select 1 from products where id=?)`, id).Scan(&exists)
	if !exists {
		return updatedProduct, fmt.Errorf("no  product found with given id")
	}
	var brand response.Brand

	selectQuery := `SELECT b.*, c.category_name
	                FROM brands b
	                JOIN categories c ON b.category_id = c.id
	                WHERE b.id = $1`

	err := c.DB.Raw(selectQuery, product.Brand).Scan(&brand).Error
	if err != nil {
		return updatedProduct, err
	}

	updateQuery := `UPDATE products SET product_name=$1,description=$2,brand=$3,category_id=$4,updated_at=NOW() WHERE id=$5
	               RETURNING id,product_name AS name,description,brand,category_id`
	err = c.DB.Raw(updateQuery, product.Name, product.Description, brand.Name, brand.Category_id, id).Scan(&updatedProduct).Error
	if err != nil {
		return updatedProduct, err
	}
	updatedProduct.CategoryName = brand.Category_name
	return updatedProduct, err
}

// DeleteProduct implements interfaces.ProductRepository.
func (c *ProductDatabase) DeleteProduct(id int) error {
	// Check if the product exists
	var exists bool
	c.DB.Raw(`SELECT EXISTS(SELECT 1 FROM products WHERE id=?)`, id).Scan(&exists)
	if !exists {
		return fmt.Errorf("no product found with given id")
	}

	// Start a transaction
	tx := c.DB.Begin()

	// Check for errors
	if tx.Error != nil {
		return tx.Error
	}

	// Delete images related to the product
	if err := tx.Exec(`DELETE FROM images WHERE product_item_id IN (SELECT id FROM product_items WHERE product_id = ?)`, id).Error; err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}

	// Delete product items related to the product
	if err := tx.Exec(`DELETE FROM product_items WHERE product_id = ?`, id).Error; err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}

	// Delete the product itself
	if err := tx.Exec(`DELETE FROM products WHERE id = ?`, id).Error; err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	return tx.Commit().Error
}

// ListAllProducts implements interfaces.ProductRepository.
func (c *ProductDatabase) ListAllProducts(queryParams helperStruct.QueryParams) ([]response.Product, error) {
	var products []response.Product
	getProductDetails := `SELECT products.product_name AS name,products.description,products.id,brand, categories.category_name
	FROM products
	JOIN categories ON products.category_id = categories.id`
	if queryParams.Query != "" && queryParams.Filter != "" {
		getProductDetails = fmt.Sprintf("%s WHERE LOWER(%s) LIKE '%%%s%%'", getProductDetails, queryParams.Filter, strings.ToLower(queryParams.Query))
	}
	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			getProductDetails = fmt.Sprintf("%s ORDER BY %s DESC", getProductDetails, queryParams.SortBy)
		} else {
			getProductDetails = fmt.Sprintf("%s ORDER BY %s ASC", getProductDetails, queryParams.SortBy)
		}
	} else {
		getProductDetails = fmt.Sprintf("%s ORDER BY products.created_at DESC", getProductDetails)
	}
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		getProductDetails = fmt.Sprintf("%s LIMIT %d OFFSET %d", getProductDetails, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		getProductDetails = fmt.Sprintf("%s LIMIT 10 OFFSET 0", getProductDetails)
	}
	err := c.DB.Raw(getProductDetails).Scan(&products).Error

	return products, err
}

// DisplayProduct implements interfaces.ProductRepository.
func (c *ProductDatabase) DisplayProduct(id int) (response.Product, error) {
	var product response.Product
	var exists bool
	c.DB.Raw(`select exists(select 1 from products where id=?)`, id).Scan(&exists)
	if !exists {
		return product, fmt.Errorf("no product found with given id")
	}
	err := c.DB.Raw(`SELECT products.product_name AS name,products.description,products.id,brand, categories.category_name
	                FROM products
	                JOIN categories ON products.category_id = categories.id
	                WHERE products.id = ?
	`, id).Scan(&product).Error
	return product, err
}

// AddProductItem implements interfaces.ProductRepository.
func (c *ProductDatabase) AddProductItem(productItem helperStruct.ProductItem) (response.ProductItem, error) {
	var newProductItem response.ProductItem
	if productItem.Price < 0 {
		return newProductItem, fmt.Errorf("price can't have a negative value")
	}
	if productItem.Ram < 0 {
		return newProductItem, fmt.Errorf("ram can't have a negative value")
	}
	if productItem.Qty < 0 {
		return newProductItem, fmt.Errorf("the quantity can't have a negative quantity")
	}
	if productItem.Storage < 0 {
		return newProductItem, fmt.Errorf("storage can't have a negative value")
	}
	if productItem.Battery < 0 {
		return newProductItem, fmt.Errorf("battery can't have a negative value")
	}
	var exists bool
	c.DB.Raw(`SELECT EXISTS(SELECT 1 FROM product_items WHERE product_id=?)`, productItem.Product_id).Scan(&exists)
	if exists {
		return newProductItem, fmt.Errorf("product_item already exists")
	}
	insertQuery := `INSERT INTO product_items (id,product_id,sku,qty_in_stock,color,ram,battery,screen_size,storage,price,graphic_processor,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW()) 
	RETURNING id,sku,color,qty_in_stock,battery,ram,screen_size,storage,price,graphic_processor`
	err := c.DB.Raw(insertQuery, productItem.Product_id, productItem.Product_id, productItem.Sku, productItem.Qty, productItem.Color, productItem.Ram, productItem.Battery, productItem.Screen_size, productItem.Storage, productItem.Price, productItem.Graphic_Processor).Scan(&newProductItem).Error
	if err != nil {
		return newProductItem, err
	}
	err = c.DB.Raw(`
    SELECT products.*, categories.category_name
    FROM products
    JOIN categories ON products.category_id = categories.id
    WHERE products.id = ?
`, productItem.Product_id).Scan(&newProductItem).Error

	return newProductItem, err
}

// UpdateProductItem implements interfaces.ProductRepository.
func (c *ProductDatabase) UpdateProductItem(id int, productItem helperStruct.ProductItem) (response.ProductItem, error) {
	var updatedProductItem response.ProductItem
	var exists bool
	c.DB.Raw(`select exists(select 1 from product_items where id=?)`, id).Scan(&exists)
	if !exists {
		return updatedProductItem, fmt.Errorf("no productitem found with given id")
	}
	updateQuery := `UPDATE product_items SET id=$1,product_id=$2,sku=$3,qty_in_stock=$4,color=$5,ram=$6,battery=$7,screen_size=$8,storage=$9,price=$10,image=$11,graphic_processor=$12 WHERE id=$13
	RETURNING id,sku,color,qty_in_stock,battery,ram,screen_size,price,image,graphic_processor,storage`
	err := c.DB.Raw(updateQuery, productItem.Product_id, productItem.Product_id, productItem.Sku, productItem.Qty, productItem.Color, productItem.Ram, productItem.Battery, productItem.Screen_size, productItem.Storage, productItem.Price, productItem.Image, productItem.Graphic_Processor, id).Scan(&updatedProductItem).Error
	if err != nil {
		return updatedProductItem, err
	}
	err = c.DB.Raw(`
    SELECT products.*, categories.category_name
    FROM products
    JOIN categories ON products.category_id = categories.id
    WHERE products.id = ?
`, productItem.Product_id).Scan(&updatedProductItem).Error
	return updatedProductItem, err
}

// ListAllProductItems implements interfaces.ProductRepository.
func (c *ProductDatabase) ListAllProductItems(queryParams helperStruct.QueryParams) ([]response.ProductItem, error) {
	var productItems []response.ProductItem
	getProductItemDetails := `
    SELECT product_items.*, products.description,products.product_name,products.brand,image_items.image,categories.category_name
    FROM product_items
    JOIN products ON product_items.product_id = products.id
    JOIN categories ON products.category_id = categories.id
	LEFT JOIN image_items ON product_items.id=image_items.product_item_id AND image_items.is_default=true
`
	if queryParams.Query != "" && queryParams.Filter != "" {
		getProductItemDetails = fmt.Sprintf("%s WHERE LOWER(%s) LIKE '%%%s%%'", getProductItemDetails, queryParams.Filter, strings.ToLower(queryParams.Query))
	}
	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			getProductItemDetails = fmt.Sprintf("%s ORDER BY %s DESC", getProductItemDetails, queryParams.SortBy)
		} else {
			getProductItemDetails = fmt.Sprintf("%s ORDER BY %s ASC", getProductItemDetails, queryParams.SortBy)
		}
	} else {
		getProductItemDetails = fmt.Sprintf("%s ORDER BY product_items.created_at DESC", getProductItemDetails)
	}
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		getProductItemDetails = fmt.Sprintf("%s LIMIT %d OFFSET %d", getProductItemDetails, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		getProductItemDetails = fmt.Sprintf("%s LIMIT 10 OFFSET 0", getProductItemDetails)
	}
	err := c.DB.Raw(getProductItemDetails).Scan(&productItems).Error
	return productItems, err

}

// DeleteProductItem implements interfaces.ProductRepository.
func (c *ProductDatabase) DeleteProductItem(id int) error {
	// Check if the product item exists
	var exists bool
	c.DB.Raw(`SELECT EXISTS(SELECT 1 FROM product_items WHERE id=?)`, id).Scan(&exists)
	if !exists {
		return fmt.Errorf("no product item found with given id")
	}

	// Start a transaction
	tx := c.DB.Begin()

	// Check for errors
	if tx.Error != nil {
		return tx.Error
	}

	// Delete images related to the product item
	if err := tx.Exec(`DELETE FROM images WHERE product_item_id = ?`, id).Error; err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}

	// Delete the product item itself
	if err := tx.Exec(`DELETE FROM product_items WHERE id = ?`, id).Error; err != nil {
		// Rollback the transaction in case of an error
		tx.Rollback()
		return err
	}

	// Commit the transaction if everything is successful
	return tx.Commit().Error
}

// DisplayProductItem implements interfaces.ProductRepository.
func (c *ProductDatabase) DisplayProductItem(id int) (response.DisplayProductItem, error) {
	var productItem response.ProductItem
	var exists bool
	c.DB.Raw(`select exists(select 1 from product_items where id=?)`, id).Scan(&exists)
	if !exists {
		return response.DisplayProductItem{}, fmt.Errorf("no productitem found with given id")
	}
	selectQuery := `
    SELECT product_items.*, products.description,products.product_name,products.brand,image_items.image,categories.category_name
    FROM product_items
    JOIN products ON product_items.product_id = products.id
    JOIN categories ON products.category_id = categories.id
	LEFT JOIN image_items ON product_items.id=image_items.product_item_id AND is_default=true
	WHERE product_items.id=?
`
	err := c.DB.Raw(selectQuery, id).Scan(&productItem).Error
	if err != nil {
		return response.DisplayProductItem{}, err
	}
	var images []response.Image
	displayImages := `SELECT  id,image FROM image_items WHERE product_item_id=?`
	err = c.DB.Raw(displayImages, id).Scan(&images).Error
	var responseProduct response.DisplayProductItem
	responseProduct.Images = images
	responseProduct.ProductSpecs = productItem
	return responseProduct, err
}

// // UploadImage implements interfaces.ProductRepository.
// func (c *ProductDatabase) UploadImage(Image helperStruct.ImageHelper) (response.ImageResponse, error) {
// 	var image response.ImageResponse
// 	err := c.DB.Raw(`INSERT INTO images(product_item_id,image) VALUES ($1,$2) RETURNING image,product_item_id AS id`, Image.ProductItemId, Image.ImageData).Scan(&image).Error
// 	return image, err
// }
// -------------------------- Uploaded-Model --------------------------//

func (c *ProductDatabase) UploadImage(filepath string, productId int) (response.Image, error) {
	updateImageDefault := `UPDATE image_items SET is_default=false WHERE product_item_id=?`
	c.DB.Exec(updateImageDefault, productId)
	var image response.Image
	uploadImage := `INSERT INTO image_items (product_item_id,image,is_default)VALUES($1,$2,true) RETURNING product_item_id AS id,image`
	err := c.DB.Raw(uploadImage, productId, filepath).Scan(&image).Error
	return image, err
}

// DeleteImage implements interfaces.ProductRepository.
func (c *ProductDatabase) DeleteImage(id int) error {
	var exists bool

	c.DB.Raw(`SELECT exists (select  1 from image_items WHERE id=?)`, id).Scan(&exists)
	if !exists {
		return fmt.Errorf("no image found for the given productitem")
	}

	err := c.DB.Exec(`DELETE FROM image_items WHERE id=?`, id).Error
	return err
}

// SearchProducts implements interfaces.ProductRepository.
func (c *ProductDatabase) SearchProducts(queryParams helperStruct.QueryParams, searchProducts string) ([]response.Product, error) {
	var products []response.Product
	search := "%" + searchProducts + "%"
	getProductDetails := fmt.Sprintf(`SELECT products.product_name AS name,products.description,products.id,brand, categories.category_name
	FROM products
	JOIN categories ON products.category_id = categories.id
	WHERE products.product_name ILIKE '%s'`, search)
	if queryParams.Query != "" && queryParams.Filter != "" {
		getProductDetails = fmt.Sprintf("%s AND LOWER(%s) LIKE '%%%s%%'", getProductDetails, queryParams.Filter, strings.ToLower(queryParams.Query))
	}
	if queryParams.SortBy != "" {
		if queryParams.SortDesc {
			getProductDetails = fmt.Sprintf("%s ORDER BY %s DESC", getProductDetails, queryParams.SortBy)
		} else {
			getProductDetails = fmt.Sprintf("%s ORDER BY %s ASC", getProductDetails, queryParams.SortBy)
		}
	} else {
		getProductDetails = fmt.Sprintf("%s ORDER BY products.created_at DESC", getProductDetails)
	}
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		getProductDetails = fmt.Sprintf("%s LIMIT %d OFFSET %d", getProductDetails, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		getProductDetails = fmt.Sprintf("%s LIMIT 10 OFFSET 0", getProductDetails)
	}
	err := c.DB.Raw(getProductDetails).Scan(&products).Error
	if len(products) == 0 {
		return []response.Product{}, fmt.Errorf("no product found with the given name")
	}

	return products, err
}
