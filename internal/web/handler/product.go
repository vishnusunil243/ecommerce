package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
)

type ProductHandler struct {
	productUseCase services.ProductUsecase
}

func NewProductHandler(productUseCase services.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
	}
}
func (cr *ProductHandler) CreateCategory(c *gin.Context) {
	var Category helperStruct.Category
	err := c.BindJSON(&Category)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "unable to bind json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	NewCategory, err := cr.productUseCase.CreateCategory(Category)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Can't create category",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, response.Response{
		StatusCode: 200,
		Message:    "category added successfully",
		Data:       NewCategory,
		Errors:     nil,
	})
}
func (cr *ProductHandler) UpdateCategory(c *gin.Context) {
	var category helperStruct.Category
	err := c.BindJSON(&category)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "bind failed",
			Data:       nil,
			Errors:     err.Error(),
		})
	}
	updatedCategory, err := cr.productUseCase.UpdateCategory(category, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating category",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "category updated successfully",
		Data:       updatedCategory,
		Errors:     nil,
	})
}
func (cr *ProductHandler) DeleteCategory(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error getting params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	err = cr.productUseCase.DeleteCategory(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error deleting category",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "category deleted successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (cr *ProductHandler) ListAllCategories(c *gin.Context) {
	categories, err := cr.productUseCase.ListAllCategories()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "unable to list all categories",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "Categories listed successfully",
		Data:       categories,
		Errors:     nil,
	})
}
func (cr *ProductHandler) DisplayCategory(c *gin.Context) {
	paramId := c.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error getting params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	category, err := cr.productUseCase.DisplayCategory(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying category information",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "category information displayed successfully",
		Data:       category,
		Errors:     nil,
	})
}
func (cr *ProductHandler) CreateBrand(c *gin.Context) {
	var Brand helperStruct.Brand
	err := c.BindJSON(&Brand)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	newBrand, err := cr.productUseCase.CreateBrand(Brand)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error creating brand",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, response.Response{
		StatusCode: 201,
		Message:    "brand created successfully",
		Data:       newBrand,
		Errors:     nil,
	})

}
func (cr *ProductHandler) UpdateBrand(c *gin.Context) {
	paramId := c.Param("id")
	Id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing parameter",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	var brand helperStruct.Brand
	err = c.BindJSON(&brand)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	updatedBrand, err := cr.productUseCase.UpdatedBrand(brand, Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating brand",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "brand updated successfully",
		Data:       updatedBrand,
		Errors:     nil,
	})
}
func (cr *ProductHandler) DeleteBrand(c *gin.Context) {
	paramId := c.Param("id")
	Id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = cr.productUseCase.DeleteBrand(Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error deleting brand",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "brand deleted successfully",
		Data:       nil,
		Errors:     nil,
	})

}
func (cr *ProductHandler) ListAllBrands(c *gin.Context) {
	var queryParams helperStruct.QueryParams
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
	allBrands, totalCount, err := cr.productUseCase.ListAllBrands(queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing all brands",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}
	responseStruct := struct {
		Brands    []response.Brand
		NoOfPages int
	}{
		Brands:    allBrands,
		NoOfPages: totalCount / queryParams.Limit,
	}
	if responseStruct.NoOfPages == 0 {
		responseStruct.NoOfPages = 1
	} else if totalCount%queryParams.Limit != 0 {
		responseStruct.NoOfPages = responseStruct.NoOfPages + 1
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "all brands listed successfully",
		Data:       responseStruct,
		Errors:     nil,
	})
}
func (p *ProductHandler) DisplayBrand(c *gin.Context) {
	paramId := c.Param("brand_id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	brand, err := p.productUseCase.DisplayBrand(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error loading brand",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "brand information loaded successfully",
		Data:       brand,
		Errors:     nil,
	})
}
func (p *ProductHandler) AddProduct(c *gin.Context) {
	var product helperStruct.Product
	err := c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	newProduct, err := p.productUseCase.AddProduct(product)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error adding product",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "product added successfully",
		Data:       newProduct,
		Errors:     nil,
	})
}
func (p *ProductHandler) UpdateProduct(c *gin.Context) {
	paramId := c.Param("product_id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	var product helperStruct.Product
	err = c.BindJSON(&product)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	updatedProduct, err := p.productUseCase.UpdateProduct(product, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating product",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "product updated successfully",
		Data:       updatedProduct,
		Errors:     nil,
	})
}
func (p *ProductHandler) DeleteProduct(c *gin.Context) {
	paramId := c.Param("product_id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = p.productUseCase.DeletProduct(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error deleting product",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "product deleted successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (p *ProductHandler) ListAllProducts(c *gin.Context) {
	var queryParams helperStruct.QueryParams
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
	queryParams.SortBy = c.Query("sort_by")
	queryParams.Query = c.Query("query")
	queryParams.Filter = c.Query("filter")
	if c.Query("sort_desc") != "" {
		queryParams.SortDesc = true
	}
	products, totalCount, err := p.productUseCase.ListAllProducts(queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing all products",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}
	responseStruct := struct {
		Products  []response.Product
		NoOfPages int
	}{
		Products:  products,
		NoOfPages: totalCount / queryParams.Limit,
	}
	if responseStruct.NoOfPages == 0 {
		responseStruct.NoOfPages = 1
	} else if totalCount%queryParams.Limit != 0 {
		responseStruct.NoOfPages = responseStruct.NoOfPages + 1
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "products listed successfully",
		Data:       responseStruct,
		Errors:     nil,
	})
}
func (p *ProductHandler) DisplayProduct(c *gin.Context) {
	paramId := c.Param("product_id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	product, err := p.productUseCase.DisplayProduct(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying product info",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "product displayed successfully",
		Data:       product,
		Errors:     nil,
	})
}
func (p *ProductHandler) AddProductItem(c *gin.Context) {
	var productItem helperStruct.ProductItem
	err := c.BindJSON(&productItem)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	newProductItem, err := p.productUseCase.AddProductItem(productItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error adding new productItem",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "new productItem added successfully",
		Data:       newProductItem,
		Errors:     nil,
	})
}
func (p *ProductHandler) UpdateProductItem(c *gin.Context) {
	paramId := c.Param("productItem_id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	var productItem helperStruct.ProductItem
	err = c.BindJSON(&productItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	updatedProductItem, err := p.productUseCase.UpdateProductItem(id, productItem)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating productItem",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "productItem udpated successfully",
		Data:       updatedProductItem,
		Errors:     nil,
	})
}
func (p *ProductHandler) ListAllProductItems(c *gin.Context) {
	var queryParams helperStruct.QueryParams
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
	queryParams.SortBy = c.Query("sort_by")
	if c.Query("sort_desc") != "" {
		queryParams.SortDesc = true
	}

	productItems, totalCount, err := p.productUseCase.ListAllProductItems(queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing all productitems",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}
	responseStruct := struct {
		ProductItems []response.ProductItem
		NoOfPages    int
	}{
		ProductItems: productItems,
		NoOfPages:    totalCount / queryParams.Limit,
	}
	if responseStruct.NoOfPages == 0 {
		responseStruct.NoOfPages = 1
	} else if totalCount%queryParams.Limit != 0 {
		responseStruct.NoOfPages = responseStruct.NoOfPages + 1
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "productitems listed successfully",
		Data:       responseStruct,
		Errors:     nil,
	})
}
func (p *ProductHandler) DeleteProductItem(c *gin.Context) {
	paramId := c.Param("productItem_id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	err = p.productUseCase.DeleteProductItem(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error deleting productitem",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "productitem deleted successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (p *ProductHandler) DisplayProductItem(c *gin.Context) {
	paramId := c.Param("productItem_id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	productItem, err := p.productUseCase.DisplayProductItem(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying productitem",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "productitem displayed successfully",
		Data:       productItem,
		Errors:     nil,
	})
}

// -------------------------- Upload-Image --------------------------//

func (cr *ProductHandler) UploadImage(c *gin.Context) {

	id := c.Param("productItem_id")
	productId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "cant find product id",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	// // Initialize MinIO client object
	// minioClient, err := minio.New(viper.GetString("ENDPOINT"), &minio.Options{
	// 	Creds:  credentials.NewStaticV4(viper.GetString("ACCESSKEY"), viper.GetString("SECRETKEY"), ""),
	// 	Secure: false, // Change to true if using HTTPS
	// })
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, response.Response{
	// 		StatusCode: 500,
	// 		Message:    "failed to initialize MinIO client",
	// 		Data:       nil,
	// 		Errors:     err.Error(),
	// 	})
	// 	return
	// }
	// driveService, err := drive.NewService(context.Background())
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, response.Response{
	// 		StatusCode: 500,
	// 		Message:    "failed to initialize Google Drive client",
	// 		Data:       nil,
	// 		Errors:     err.Error(),
	// 	})
	// 	return
	// }

	var Image response.Image

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "initialization error",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	files, ok := form.File["images"]
	if !ok || len(files) == 0 {
		// Handle the case where "images" key is not found or no files are present
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "No files found in the 'images' key",
			Data:       nil,
			Errors:     "No files found",
		})
		return
	}
	images := make([]string, 0)
	for _, file := range files {
		// Upload the file to specific dst.
		filePath := "../asset/uploads/" + file.Filename
		c.SaveUploadedFile(file, filePath)

		// objectName := fmt.Sprintf("prroduct_%d_%s", productId, file.Filename)
		// Upload file to MinIO

		// Open the file to get an io.Reader
		// fileData, err := file.Open()
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, response.Response{
		// 		StatusCode: 400,
		// 		Message:    "can't open form file",
		// 		Data:       nil,
		// 		Errors:     err.Error(),
		// 	})
		// 	return
		// }
		// defer fileData.Close()
		// Set the content type
		// contentType := "image/webp"
		// _, err = minioClient.PutObject(context.TODO(), viper.GetString("BUCKETNAME"), objectName, fileData, file.Size, minio.PutObjectOptions{
		// 	ContentType: contentType,
		// })
		// // Create a new Google Drive file
		// file1 := &drive.File{
		// 	Name: file.Filename,
		// }

		// res, err := driveService.Files.Create(file1).Media(fileData).Do()

		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, response.Response{
		// 		StatusCode: 400,
		// 		Message:    "can't upload images to google drive",
		// 		Data:       nil,
		// 		Errors:     err.Error(),
		// 	})
		// 	return
		// }
		// if err != nil {
		// 	c.JSON(http.StatusBadRequest, response.Response{
		// 		StatusCode: 400,
		// 		Message:    "can't upload images to MinIO",
		// 		Data:       nil,
		// 		Errors:     err.Error(),
		// 	})
		// 	return
		// }

		// Get the MinIO URL for the uploaded file
		// objectURL := fmt.Sprintf("%s/%s/%s", viper.GetString("ENDPOINT"), viper.GetString("BUCKETNAME"), objectName)
		// objectURL := res.WebViewLink

		Image, err = cr.productUseCase.UploadImage(filePath, productId)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				StatusCode: 400,
				Message:    "cant upload images",
				Data:       nil,
				Errors:     err.Error(),
			})
			return
		}
		images = append(images, Image.Image)

	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "image uploaded",
		Data:       images,
		Errors:     nil,
	})

}
func (p *ProductHandler) DeleteImage(c *gin.Context) {
	paramId := c.Param("image_id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = p.productUseCase.DeleteImage(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error deleting image",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "image deleted successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (p *ProductHandler) SearchProducts(c *gin.Context) {
	var queryParams helperStruct.QueryParams
	var search helperStruct.SearchProducts
	err := c.BindJSON(&search)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
	queryParams.SortBy = c.Query("sort_by")
	if c.Query("sort_desc") != "" {
		queryParams.SortDesc = true
	}
	queryParams.Query = c.Query("query")
	queryParams.Filter = c.Query("filter")
	searchProducts, err := p.productUseCase.SearchProducts(queryParams, search.SearchProducts)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error searching products",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "products fetched successfully",
		Data:       searchProducts,
		Errors:     nil,
	})
}
