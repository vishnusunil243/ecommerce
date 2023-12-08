package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
	"main.go/internal/web/handlerUtil"
)

type WishlistHandler struct {
	wishlistUseCase services.WishlistUseCase
}

func NewWishlistHandler(wishListUseCase services.WishlistUseCase) *WishlistHandler {
	return &WishlistHandler{
		wishlistUseCase: wishListUseCase,
	}
}
func (w *WishlistHandler) AddToWishlist(c *gin.Context) {
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
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving userId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = w.wishlistUseCase.AddToWishlist(productId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error adding to wishlist",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "product added to wishlist",
		Data:       nil,
		Errors:     nil,
	})
}
func (w *WishlistHandler) RemoveFromWishlist(c *gin.Context) {
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
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving userId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	err = w.wishlistUseCase.RemoveFromWishlist(productId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error removing from wishlist",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "product removed from wishlist",
		Data:       nil,
		Errors:     nil,
	})
}
func (w *WishlistHandler) ListAllWishlist(c *gin.Context) {
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving userId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	wishlists, err := w.wishlistUseCase.ListAllWishlist(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying wishlists",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if len(wishlists) == 0 {
		c.JSON(http.StatusOK, response.Response{
			StatusCode: 200,
			Message:    "there are no items in wishlist",
			Data:       nil,
			Errors:     nil,
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "wishlist fetched successfully",
		Data:       wishlists,
		Errors:     nil,
	})
}
func (w *WishlistHandler) DisplayWishlistProduct(c *gin.Context) {
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
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving userId",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	wishlist, err := w.wishlistUseCase.DisplayWishlistProduct(productId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying wishlist product",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "wishlist product fetched successfully",
		Data:       wishlist,
		Errors:     nil,
	})
}
