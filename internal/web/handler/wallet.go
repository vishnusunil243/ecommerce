package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
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
	var queryParams helperStruct.QueryParams
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
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
	walletHistory, totalCount, err := w.walletUseCase.WalletHistory(userId, queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying wallet history",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}
	responseStruct := struct {
		WalletHistories []response.WalletHistories
		NoOfPages       int
	}{
		WalletHistories: walletHistory,
		NoOfPages:       totalCount / queryParams.Limit,
	}
	if responseStruct.NoOfPages == 0 {
		responseStruct.NoOfPages = 1
	} else if totalCount%queryParams.Limit != 0 {
		responseStruct.NoOfPages = responseStruct.NoOfPages + 1
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "wallet history fetched successfully",
		Data:       responseStruct,
		Errors:     nil,
	})
}
