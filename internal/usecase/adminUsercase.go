package usecase

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/infrastructure/config"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type adminUsecase struct {
	adminRepo interfaces.AdminRepository
}

func NewAdminUsecase(adminRepo interfaces.AdminRepository) services.AdminUseCase {
	return &adminUsecase{
		adminRepo: adminRepo,
	}
}

// AdminLogin implements interfaces.AdminUseCase.
func (c *adminUsecase) AdminLogin(admin helperStruct.LoginReq) (string, error) {
	var cfg config.Config
	adminData, err := c.adminRepo.AdminLogin(admin.Email)
	if err != nil {
		return "", err
	}
	if adminData.IsBlocked {
		return "", fmt.Errorf("admin is blocked")
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminData.Password), []byte(admin.Password))
	fmt.Println(adminData.Password)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"id":  adminData.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.SECRET))
	fmt.Println("sooo")
	if err != nil {
		return "", err
	}

	return ss, nil

}

// ListAllUsers implements interfaces.AdminUseCase.
func (cr *adminUsecase) ListAllUsers(queryParams helperStruct.QueryParams) ([]response.UserDetails, error) {
	users, err := cr.adminRepo.ListAllUsers(queryParams)
	return users, err
}

// DisplayUser implements interfaces.AdminUseCase.
func (cr *adminUsecase) DisplayUser(id int) (response.UserDetails, error) {
	user, err := cr.adminRepo.DispalyUser(id)
	return user, err
}

// ReportUser implements interfaces.AdminUseCase.
func (cr *adminUsecase) ReportUser(UsersId int) (response.UserReport, error) {
	reportInfo, err := cr.adminRepo.ReportUser(UsersId)
	return reportInfo, err
}

// GetDashboard implements interfaces.AdminUseCase.
func (cr *adminUsecase) GetDashboard(dashboard helperStruct.Dashboard) (response.DashBoard, error) {
	newDashboard, err := cr.adminRepo.GetDashBoard(dashboard)
	return newDashboard, err
}

// ViewSalesReport implements interfaces.AdminUseCase.
func (cr *adminUsecase) ViewSalesReport(filter helperStruct.Dashboard) ([]response.SalesReport, error) {
	salesReports, err := cr.adminRepo.ViewSalesReport(filter)
	return salesReports, err
}
