package repository

import (
	"fmt"

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
	var id int
	selectQuery := `SELECT id FROM categories WHERE category_name=$1`
	err := c.DB.Raw(selectQuery, brand.Category_name).Scan(&id).Error
	if err != nil {
		return newbrand, err
	}
	insertQuery := `INSERT INTO brands (brandname,description,category_id,created_at) VALUES ($1,$2,$3,NOW()) RETURNING id,brandname AS name,description,category_id`
	err = c.DB.Raw(insertQuery, brand.Name, brand.Description, id).Scan(&newbrand).Error
	if err != nil {
		return newbrand, err
	}
	newbrand.Category_name = brand.Category_name
	return newbrand, nil

}

// UpdateBrand implements interfaces.ProductRepository.
func (c *ProductDatabase) UpdateBrand(brand helperStruct.Brand, id int) (response.Brand, error) {
	var updatedBrand response.Brand
	var categoryid int
	selectQuery := `SELECT id FROM categories WHERE category_name=$1`
	err := c.DB.Raw(selectQuery, brand.Category_name).Scan(&categoryid).Error
	if err != nil {
		return updatedBrand, err
	}
	updateQuery := `UPDATE brands SET brandname=$1,description=$2,category_id=$3,updated_at=NOW() WHERE id=$4 RETURNING id,brandname AS name,category_id,description`
	err = c.DB.Raw(updateQuery, brand.Name, brand.Description, categoryid, id).Scan(&updatedBrand).Error
	updatedBrand.Category_name = brand.Category_name
	if updatedBrand.Id == 0 {
		return updatedBrand, fmt.Errorf("no such brand to update")
	}
	return updatedBrand, err
}

// DeleteBrands implements interfaces.ProductRepository.
func (c *ProductDatabase) DeleteBrand(id int) error {
	var exists bool
	c.DB.Raw(`select exists(select 1 from categories where id=?)`, id).Scan(&exists)
	if !exists {
		return fmt.Errorf("no brand found with given id")
	}
	err := c.DB.Exec(`DELETE FROM brands WHERE id=?`, id).Error
	return err
}

// ListAllBrands implements interfaces.ProductRepository.
func (c *ProductDatabase) ListAllBrands() ([]response.Brand, error) {
	var brands []response.Brand
	err := c.DB.Raw(`
    SELECT brands.brandname AS name,brands.id,brands.category_id,brands.description, categories.category_name
    FROM brands
    JOIN categories ON brands.category_id = categories.id
	
`).Scan(&brands).Error
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
	err := c.DB.Raw(`
    SELECT b.*, c.category_name
    FROM brands b
    JOIN categories c ON b.category_id = c.id
    WHERE b.brandname = ?
`, product.Brand).Scan(&brand).Error
	if err != nil {
		return newProduct, err
	}
	insertQuery := `INSERT INTO products (product_name,description,brand,category_id,created_at) VALUES ($1,$2,$3,$4,NOW())
	RETURNING id,product_name AS name,description,brand,category_id`
	err = c.DB.Raw(insertQuery, product.Name, product.Description, product.Brand, brand.Category_id).Scan(&newProduct).Error
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
	                WHERE b.brandname = $1`

	err := c.DB.Raw(selectQuery, product.Brand).Scan(&brand).Error
	if err != nil {
		return updatedProduct, err
	}

	updateQuery := `UPDATE products SET product_name=$1,description=$2,brand=$3,category_id=$4,updated_at=NOW() WHERE id=$5
	               RETURNING id,product_name AS name,description,brand,category_id`
	err = c.DB.Raw(updateQuery, product.Name, product.Description, product.Brand, brand.Category_id, id).Scan(&updatedProduct).Error
	if err != nil {
		return updatedProduct, err
	}
	updatedProduct.CategoryName = brand.Category_name
	return updatedProduct, err
}

// DeleteProduct implements interfaces.ProductRepository.
func (c *ProductDatabase) DeleteProduct(id int) error {
	var exists bool
	c.DB.Raw(`select exists(select 1 from products where id=?)`, id).Scan(&exists)
	if !exists {
		return fmt.Errorf("no product found with given id")
	}
	err := c.DB.Exec(`DELETE FROM products WHERE id= ?`, id).Error
	return err
}

// ListAllProducts implements interfaces.ProductRepository.
func (c *ProductDatabase) ListAllProducts() ([]response.Product, error) {
	var products []response.Product
	err := c.DB.Raw(`SELECT products.product_name AS name,products.description,products.id,brand, categories.category_name
	FROM products
	JOIN categories ON products.category_id = categories.id`).Scan(&products).Error
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
	insertQuery := `INSERT INTO product_items (id,product_id,sku,qty_in_stock,color,ram,battery,screen_size,storage,price,image,graphic_processor,created_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,NOW()) 
	RETURNING id,sku,color,qty_in_stock,battery,ram,screen_size,storage,price,image,graphic_processor`
	err := c.DB.Raw(insertQuery, productItem.Product_id, productItem.Product_id, productItem.Sku, productItem.Qty, productItem.Color, productItem.Ram, productItem.Battery, productItem.Screen_size, productItem.Storage, productItem.Price, productItem.Image, productItem.Graphic_Processor).Scan(&newProductItem).Error
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
func (c *ProductDatabase) ListAllProductItems() ([]response.ProductItem, error) {
	var productItems []response.ProductItem
	selectQuery := `
    SELECT product_items.*, products.description,products.product_name,products.brand, categories.category_name
    FROM product_items
    JOIN products ON product_items.product_id = products.id
    JOIN categories ON products.category_id = categories.id
`
	err := c.DB.Raw(selectQuery).Scan(&productItems).Error
	return productItems, err

}

// DeleteProductItem implements interfaces.ProductRepository.
func (c *ProductDatabase) DeleteProductItem(id int) error {
	var exists bool
	c.DB.Raw(`select exists(select 1 from product_items where id=?)`, id).Scan(&exists)
	if !exists {
		return fmt.Errorf("no productitem found with given id")
	}
	err := c.DB.Exec(`DELETE FROM product_items WHERE id=?`, id).Error
	return err
}

// DisplayProductItem implements interfaces.ProductRepository.
func (c *ProductDatabase) DisplayProductItem(id int) (response.ProductItem, error) {
	var productItem response.ProductItem
	var exists bool
	c.DB.Raw(`select exists(select 1 from product_items where id=?)`, id).Scan(&exists)
	if !exists {
		return productItem, fmt.Errorf("no productitem found with given id")
	}
	selectQuery := `
    SELECT product_items.*, products.description,products.product_name,products.brand, categories.category_name
    FROM product_items
    JOIN products ON product_items.product_id = products.id
    JOIN categories ON products.category_id = categories.id
	WHERE product_items.id=?
`
	err := c.DB.Raw(selectQuery, id).Scan(&productItem).Error
	return productItem, err
}
