package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
	"main.go/internal/web/handlerUtil"
	"main.go/internal/web/middleware"
)

type UserHandler struct {
	userUseCase   services.UserUseCase
	cartUseCase   services.CartUseCase
	walletUsecase services.WalletUseCase
}

func NewUserHandler(usecase services.UserUseCase, cartusecase services.CartUseCase, walletUsecase services.WalletUseCase) *UserHandler {
	return &UserHandler{
		userUseCase:   usecase,
		cartUseCase:   cartusecase,
		walletUsecase: walletUsecase,
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
	err = cr.cartUseCase.CreateCart(userData.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error creating cart for the user",
			Data:       nil,
			Errors:     err.Error(),
		})
		fmt.Println(err.Error())
		return
	}
	err = cr.walletUsecase.CreateWallet(userData.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error creating wallet",
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
func (u *UserHandler) AddAddress(c *gin.Context) {
	var address helperStruct.Address
	err := c.BindJSON(&address)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	Id, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error getting user id",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	newAdress, err := u.userUseCase.AddAdress(Id, address)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "errror adding address",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "address addedd successfully",
		Data:       newAdress,
		Errors:     nil,
	})

}
func (u *UserHandler) UpdateAddress(c *gin.Context) {
	var address helperStruct.Address
	err := c.BindJSON(&address)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	paramId := c.Param("address_id")
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
	updatedAddress, err := u.userUseCase.UpdateAddress(userId, id, address)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "errror updating address",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "address updated successfully",
		Data:       updatedAddress,
		Errors:     nil,
	})

}
func (u *UserHandler) DeleteAddress(c *gin.Context) {
	paramId := c.Param("address_id")
	addressId, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing address id",
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
	}
	err = u.userUseCase.DeleteAddress(addressId, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error deleting address",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "address deleted successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (u *UserHandler) ViewUserProfile(c *gin.Context) {
	Id, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving users id",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	userProfile, err := u.userUseCase.ViewUserProfile(Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving user profile",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "user profile retrieved successfully",
		Data:       userProfile,
		Errors:     nil,
	})

}
func (u *UserHandler) ListAllAddresses(c *gin.Context) {
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
	addresses, err := u.userUseCase.ListAllAddresses(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying addresses",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "addresses displayed successfully",
		Data:       addresses,
		Errors:     nil,
	})
}
func (u *UserHandler) UpdateMobile(c *gin.Context) {
	var mobile helperStruct.UpdateMobile
	err := c.BindJSON(&mobile)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	Id, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving user id",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	userProfile, err := u.userUseCase.UpdateMobile(Id, mobile.Mobile)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating mobile",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "mobile number updated successfully",
		Data:       userProfile,
		Errors:     nil,
	})
}
func (u *UserHandler) ChangePassword(c *gin.Context) {
	var password helperStruct.UpdatePassword
	err := c.BindJSON(&password)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	Id, err := handlerUtil.GetUserIdFromContext(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error getting user id from context",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	userProfile, err := u.userUseCase.ChangePassword(Id, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating password",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "password changed successfully",
		Data:       userProfile,
		Errors:     nil,
	})
}
func (u *UserHandler) ForgotPassword(c *gin.Context) {
	var forgotPassword helperStruct.ForgotPassword
	err := c.BindJSON(&forgotPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if forgotPassword.OTP == "" {
		err = middleware.SendOTP(forgotPassword.Email)
		if err != nil {
			c.JSON(http.StatusBadRequest, response.Response{
				StatusCode: 400,
				Message:    "error sending otp to the given email",
				Data:       nil,
				Errors:     err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, response.Response{
			StatusCode: 200,
			Message:    "please enter the otp",
			Data:       nil,
			Errors:     nil,
		})
		return
	} else {
		if !middleware.VerifyOTP(forgotPassword.Email, forgotPassword.OTP) {
			c.JSON(http.StatusBadRequest, response.Response{
				StatusCode: 400,
				Message:    "invalid otp",
				Data:       nil,
				Errors:     err.Error(),
			})
			return
		}
	}
	err = u.userUseCase.ForgotPassword(forgotPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error updating the password",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "password changed successfully",
		Data:       nil,
		Errors:     nil,
	})

}
