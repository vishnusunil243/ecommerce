package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
)

type DiscountHandler struct {
	discountUsecase services.DiscountUseCase
}

func NewDiscountHandler(discountUsecase services.DiscountUseCase) *DiscountHandler {
	return &DiscountHandler{
		discountUsecase: discountUsecase,
	}
}
func (d *DiscountHandler) AddDiscount(c *gin.Context) {
	var discount helperStruct.Discount
	err := c.BindJSON(&discount)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	newDiscount, err := d.discountUsecase.AddDiscount(discount)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error adding discount",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "discount added successfully",
		Data:       newDiscount,
		Errors:     nil,
	})
}
func (d *DiscountHandler) UpdateDiscount(c *gin.Context) {
	paramId := c.Param("discountId")
	discountId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	var discount helperStruct.Discount
	err = c.BindJSON(&discount)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	updatedDiscount, err := d.discountUsecase.UpdateDiscount(discount, uint(discountId))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating discount",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "discount updated successfully",
		Data:       updatedDiscount,
		Errors:     nil,
	})

}
func (d *DiscountHandler) DeleteDiscount(c *gin.Context) {
	paramId := c.Param("discountId")
	discountId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = d.discountUsecase.DeleteDiscount(discountId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error deleting discount",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "discount deleted successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (d *DiscountHandler) ListAllDiscounts(c *gin.Context) {
	var queryParams helperStruct.QueryParams
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	discounts, totalCount, err := d.discountUsecase.ListAllDiscount(queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying discounts",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}
	responseStruct := struct {
		Discounts []response.Discount
		NoOfPages int
	}{
		Discounts: discounts,
		NoOfPages: totalCount / queryParams.Limit,
	}
	if responseStruct.NoOfPages == 0 {
		responseStruct.NoOfPages = 1
	} else if totalCount%queryParams.Limit != 0 {
		responseStruct.NoOfPages = responseStruct.NoOfPages + 1
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "discounts displayed successfully",
		Data:       responseStruct,
		Errors:     nil,
	})
}
