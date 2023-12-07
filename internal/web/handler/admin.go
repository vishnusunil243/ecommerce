package handler

import (
	"encoding/csv"
	"net/http"
	"strconv"

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
func (cr *AdminHandler) AdminLogout(c *gin.Context) {
	c.SetCookie("AdminAuth", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "admin logged out successfully",
		Data:       nil,
		Errors:     nil,
	})
}
func (cr *AdminHandler) ListAllUsers(c *gin.Context) {
	var queryParams helperStruct.QueryParams
	queryParams.Page, _ = strconv.Atoi(c.Query("page"))
	queryParams.Limit, _ = strconv.Atoi(c.Query("limit"))
	users, totalCount, err := cr.adminUsecase.ListAllUsers(queryParams)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error listing users",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	if queryParams.Limit == 0 {
		queryParams.Limit = 10
	}
	responseStruct := struct {
		Users     []response.UserDetails
		NoOfPages int
	}{
		Users:     users,
		NoOfPages: totalCount / queryParams.Limit,
	}
	if responseStruct.NoOfPages == 0 {
		responseStruct.NoOfPages = 1
	} else if totalCount%queryParams.Limit != 0 {
		responseStruct.NoOfPages = responseStruct.NoOfPages + 1
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "users listed successfully",
		Data:       responseStruct,
		Errors:     nil,
	})
}
func (cr *AdminHandler) DisplayUser(c *gin.Context) {
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
	user, err := cr.adminUsecase.DisplayUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying user",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "user displayed successfully",
		Data:       user,
		Errors:     nil,
	})
}
func (a *AdminHandler) ReportUser(c *gin.Context) {
	paramUsersId := c.Param("user_id")
	UsersId, err := strconv.Atoi(paramUsersId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error parsing params",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}

	reportInfo, err := a.adminUsecase.ReportUser(UsersId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error reporting user",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "user reporteed successfully",
		Data:       reportInfo,
		Errors:     nil,
	})
}
func (a *AdminHandler) GetDashboard(c *gin.Context) {
	var dashboard helperStruct.Dashboard
	dashboard.StartDate = c.Query("start_date")
	dashboard.EndDate = c.Query("end_date")
	dashboard.Day, _ = strconv.Atoi(c.Query("day"))
	dashboard.Month, _ = strconv.Atoi(c.Query("month"))
	dashboard.Year, _ = strconv.Atoi(c.Query("year"))
	newDashboard, err := a.adminUsecase.GetDashboard(dashboard)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error displaying dashboard",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "dashboard fetched successfully",
		Data:       newDashboard,
		Errors:     nil,
	})
}
func (a *AdminHandler) ViewSalesReport(c *gin.Context) {
	var filter helperStruct.Dashboard
	filter.Year, _ = strconv.Atoi(c.Query("year"))
	filter.Day, _ = strconv.Atoi(c.Query("day"))
	filter.Month, _ = strconv.Atoi(c.Query("month"))
	filter.StartDate = c.Query("start_date")
	filter.EndDate = c.Query("end_date")
	salesReports, err := a.adminUsecase.ViewSalesReport(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "error retrieving sales reports",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, response.Response{
		StatusCode: 200,
		Message:    "salesreports fetched successfully",
		Data:       salesReports,
		Errors:     nil,
	})
}
func (a *AdminHandler) DownloadSalesReport(c *gin.Context) {
	var filter helperStruct.Dashboard
	filter.Year, _ = strconv.Atoi(c.Query("year"))
	filter.Day, _ = strconv.Atoi(c.Query("day"))
	filter.Month, _ = strconv.Atoi(c.Query("month"))
	filter.StartDate = c.Query("start_date")
	filter.EndDate = c.Query("end_date")
	sales, err := a.adminUsecase.ViewSalesReport(filter)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Response{
			StatusCode: 400,
			Message:    "cant get sales report",
			Data:       nil,
			Errors:     err.Error(),
		})
		return
	}
	// Set headers so browser will download the file
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment;filename=sales.csv")

	// Create a CSV writer using our response writer as our io.Writer
	wr := csv.NewWriter(c.Writer)

	// Write CSV header row
	headers := []string{"Name", "PaymentType", "OrderDate", "OrderTotal"}
	if err := wr.Write(headers); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	// Write data rows
	for _, sale := range sales {
		row := []string{sale.Name, sale.PaymentType, sale.OrderDate.Format("2006-01-02 15:04:05"), strconv.Itoa(sale.OrderTotal)}
		if err := wr.Write(row); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	// Flush the writer's buffer to ensure all data is written to the client
	wr.Flush()
	if err := wr.Error(); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

}
