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

type superAdminUsecase struct {
	superAdminRepo interfaces.SuperAdminRepository
}

func NewSuperAdminUsecase(superAdminRepo interfaces.SuperAdminRepository) services.SuperAdminUseCase {
	return &superAdminUsecase{
		superAdminRepo: superAdminRepo,
	}
}

// SuperLogin implements interfaces.SuperAdminUseCase.
func (cr *superAdminUsecase) SuperLogin(superadmin helperStruct.SuperLoginReq) (string, error) {
	var cfg config.Config
	superAdmin, err := cr.superAdminRepo.Login(superadmin)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(superAdmin.Password), []byte(superadmin.Password))
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}
	claims := jwt.MapClaims{
		"id":  superAdmin.Id,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.SECRET))
	if err != nil {
		return "", err
	}

	return ss, nil
}

// CreateAdmin implements interfaces.SuperAdminUseCase.
func (cr *superAdminUsecase) CreateAdmin(admin helperStruct.CreateAdmin) (response.AdminData, error) {
	var newAdmin response.AdminData
	hash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return newAdmin, fmt.Errorf("error hashing password")
	}
	admin.Password = string(hash)
	newAdmin, err = cr.superAdminRepo.CreateAdmin(admin)
	return newAdmin, err
}

// ListAllAdmins implements interfaces.SuperAdminUseCase.
func (cr *superAdminUsecase) ListAllAdmins(queryParams helperStruct.QueryParams) ([]response.AdminData, error) {
	admins, err := cr.superAdminRepo.ListAllAdmins(queryParams)
	return admins, err
}

// DisplayAdmin implements interfaces.SuperAdminUseCase.
func (cr *superAdminUsecase) DisplayAdmin(id int) (response.AdminData, error) {
	admin, err := cr.superAdminRepo.DisplayAdmin(id)
	return admin, err
}

// BlockAdmin implements interfaces.SuperAdminUseCase.
func (cr *superAdminUsecase) BlockAdmin(id int) (response.AdminData, error) {
	admin, err := cr.superAdminRepo.BlockAdmin(id)
	return admin, err
}

// BlockUser implements interfaces.SuperAdminUseCase.
func (cr *superAdminUsecase) BlockUser(id int) (response.UserData, error) {
	userData, err := cr.superAdminRepo.BlockUser(id)
	return userData, err
}

// UnBlockAdminManually implements interfaces.SuperAdminUseCase.
func (cr *superAdminUsecase) UnBlockAdminManually(id int) (response.AdminData, error) {
	adminData, err := cr.superAdminRepo.UnBlockAdminManually(id)
	return adminData, err
}

// UnBlockUserManually implements interfaces.SuperAdminUseCase.
func (cr *superAdminUsecase) UnBlockUserManually(id int) (response.UserData, error) {
	userData, err := cr.superAdminRepo.UnBlockUserManually(id)
	return userData, err
}
