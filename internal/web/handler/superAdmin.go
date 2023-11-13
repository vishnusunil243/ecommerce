package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	services "main.go/internal/usecase/interface"
)

type SuperAdminHandler struct {
	superAdminUsecase services.SuperAdminUseCase
}

func NewSuperAdminHandler(superAdminUsecase services.SuperAdminUseCase) *SuperAdminHandler {
	return &SuperAdminHandler{
		superAdminUsecase: superAdminUsecase,
	}
}
func (s *SuperAdminHandler) SuperLogin(c *gin.Context) {
	var superLoginReq helperStruct.SuperLoginReq
	err := c.BindJSON(&superLoginReq)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, response.Response{
			StatusCode: 422,
			Message:    "error binding json",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	token, err := s.superAdminUsecase.SuperLogin(superLoginReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error logging in",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("SuperAdminAuth", token, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "login successfull",
		Data:       nil,
		Errors:     nil,
	})
}
func (s *SuperAdminHandler) SuperLogout(c *gin.Context) {
	c.SetCookie("SuperAdminAuth", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "superadmin logged out successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (s *SuperAdminHandler) CreateAdmin(c *gin.Context) {
	var admin helperStruct.CreateAdmin
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
	newAdmin, err := s.superAdminUsecase.CreateAdmin(admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "errro creaeting admin",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "admin created successfully",
		Data:       newAdmin,
		Errors:     nil,
	})
}
func (s *SuperAdminHandler) ListAllAdmins(c *gin.Context) {
	var queryParams helperStruct.QueryParams
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
	admins, err := s.superAdminUsecase.ListAllAdmins(queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing all admins",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "admins listed successfully",
		Data:       admins,
		Errors:     nil,
	})
}
func (p *SuperAdminHandler) DisplayAdmin(c *gin.Context) {
	paramId := c.Param("admin_id")
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
	admin, err := p.superAdminUsecase.DisplayAdmin(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying admin",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "admin displayed successfully",
		Data:       admin,
		Errors:     nil,
	})
}
func (s *SuperAdminHandler) BlockAdmin(c *gin.Context) {
	paramId := c.Param("admin_id")
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
	admin, err := s.superAdminUsecase.BlockAdmin(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error blocking admin",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "admin blocked successfully",
		Data:       admin,
		Errors:     nil,
	})
}
func (s *SuperAdminHandler) BlockUser(c *gin.Context) {
	paramId := c.Param("user_id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "errror parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	userData, err := s.superAdminUsecase.BlockUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error blocking user",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "user blocked successfully",
		Data:       userData,
		Errors:     nil,
	})
}
func (s *SuperAdminHandler) UnBlockAdminManually(c *gin.Context) {
	paramId := c.Param("admin_id")
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
	userData, err := s.superAdminUsecase.UnBlockAdminManually(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error unblocking user",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "admin unblocked successfully",
		Data:       userData,
		Errors:     nil,
	})
}
func (s *SuperAdminHandler) UnBlockUserManually(c *gin.Context) {
	paramId := c.Param("user_id")
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
	userData, err := s.superAdminUsecase.UnBlockUserManually(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error unblocking user",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "user unblocked successfully",
		Data:       userData,
		Errors:     nil,
	})
}
