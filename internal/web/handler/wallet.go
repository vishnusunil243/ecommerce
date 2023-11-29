package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
	"main.go/internal/web/handlerUtil"
)

type WalletHandler struct {
	walletUseCase services.WalletUseCase
}

func NewWalletHandler(walletUseCase services.WalletUseCase) *WalletHandler {
	return &WalletHandler{
		walletUseCase: walletUseCase,
	}
}
func (w *WalletHandler) DisplayWallet(c *gin.Context) {
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
	wallet, err := w.walletUseCase.DisplayWallet(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying wallet of the user",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "wallet displayed successfully",
		Data:       wallet,
		Errors:     nil,
	})
}
func (w *WalletHandler) WalletHistory(c *gin.Context) {
	userId, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying wallet history",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	walletHistory, err := w.walletUseCase.WalletHistory(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying wallet history",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "wallet history fetched successfully",
		Data:       walletHistory,
		Errors:     nil,
	})
}
