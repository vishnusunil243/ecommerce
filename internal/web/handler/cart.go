package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
	"main.go/internal/web/handlerUtil"
)

type CartHandler struct {
	cartUsecase services.CartUseCase
}

func NewCartHandler(cartUsecase services.CartUseCase) *CartHandler {
	return &CartHandler{
		cartUsecase: cartUsecase,
	}
}
func (cr *CartHandler) AddToCart(c *gin.Context) {
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error getting id from context",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	paramId := c.Param("product_item_id")
	productId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = cr.cartUsecase.AddToCart(productId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error adding productitem to cart",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "product added to cart",
		Data:       nil,
		Errors:     nil,
	})

}
func (cr *CartHandler) RemoveFromCart(c *gin.Context) {
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving user id from context",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	paramId := c.Param("product_item_id")
	productId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing parameters",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = cr.cartUsecase.RemoveFromCart(productId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error removing product from cart",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "product removed from cart",
		Data:       nil,
		Errors:     nil,
	})
}
func (cr *CartHandler) ListCart(c *gin.Context) {
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "errror retrieving user id from contex",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	viewCart, err := cr.cartUsecase.ListCart(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing cart",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if len(viewCart.CartItems) == 0 {
		c.JSON(http.StatusOK, response.Response{
			StatusCode: 200,
			Message:    "there are no items in cart",
			Data:       nil,
			Errors:     nil,
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "cart displayed successfully",
		Data:       viewCart,
		Errors:     nil,
	})
}
