package usecase

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"main.go/internal/common/helperStruct"
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
	fmt.Println("ssssss")
	var cfg config.Config
	adminData, err := c.adminRepo.AdminLogin(admin.Email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(adminData.Password), []byte(admin.Password))
	fmt.Println(adminData.Password)
	fmt.Println("hii")
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
