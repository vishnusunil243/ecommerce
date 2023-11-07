package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
)

type AdminHandler struct {
	adminUsecase services.AdminUseCase
}

func NewAdminHandler(adminUsecase services.AdminUseCase) *AdminHandler {
	return &AdminHandler{
		adminUsecase: adminUsecase,
	}
}
func (cr *AdminHandler) AdminLogin(c *gin.Context) {
	var admin helperStruct.LoginReq
	err := c.BindJSON(&admin)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	token, err := cr.adminUsecase.AdminLogin(admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error generating jwt",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("AdminAuth", token, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "admin signed in successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func AdminLogout(c *gin.Context) {
	c.SetCookie("AdminAuth", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "user logged out successfully",
		Data:       nil,
		Errors:     nil,
	})
}
