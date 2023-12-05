package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
	"main.go/internal/web/handlerUtil"
)

type ReferralHandler struct {
	referralUsecase services.ReferralUseCase
}

func NewReferralHandler(referralUsecase services.ReferralUseCase) *ReferralHandler {
	return &ReferralHandler{
		referralUsecase: referralUsecase,
	}
}
func (r *ReferralHandler) ReferralOffer(c *gin.Context) {
	var referralOffer helperStruct.ReferralOffer
	err := c.BindJSON(&referralOffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
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
	err = r.referralUsecase.ReferralOffer(userId, referralOffer.ReferralId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error redeeming referral offer",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "referral offer redeemed successfully an amount of rs.20 will be deposited in the wallet shortly",
		Data:       nil,
		Errors:     nil,
	})
}
