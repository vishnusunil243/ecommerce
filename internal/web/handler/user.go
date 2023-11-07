package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
	"main.go/internal/web/middleware"
)

type UserHandler struct {
	userUseCase services.UserUseCase
}

func NewUserHandler(usecase services.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: usecase,
	}
}
func (cr *UserHandler) UserSignup(c *gin.Context) {
	var user helperStruct.UserReq
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "error binding user",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if user.OTP == "" {
		err = middleware.SendOTP(user.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				StatusCode: 400,
				Message:    "unable to send otp",
				Data:       nil,
				Errors:     err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response.Response{
			StatusCode: 200,
			Message:    "please enter your otp",
			Data:       nil,
			Errors:     nil,
		})
		return

	} else {
		if !middleware.VerifyOTP(user.Email, user.OTP) {
			c.JSON(http.StatusUnauthorized, response.Response{
				StatusCode: 401,
				Message:    "otp not verified",
				Data:       nil,
				Errors:     nil,
			})
			return
		}
	}
	userData, err := cr.userUseCase.UserSignup(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "unable to signup",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, response.Response{
		StatusCode: 201,
		Message:    "user signed up successfully",
		Data:       userData,
		Errors:     nil,
	})

}
func (cr *UserHandler) UserLogin(c *gin.Context) {
	var user helperStruct.LoginReq
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	token, err := cr.userUseCase.UserLogin(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "Login failed",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("UserAuth", token, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "login successfull",
		Data:       nil,
		Errors:     nil,
	})
}
func UserLogout(c *gin.Context) {
	c.SetCookie("UserAuth", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "user logged out successfully",
		Data:       nil,
		Errors:     nil,
	})
}
