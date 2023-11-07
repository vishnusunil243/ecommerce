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

type UserUseCase struct {
	userRepo interfaces.UserRepository
}

func NewUserUsecase(repo interfaces.UserRepository) services.UserUseCase {
	return &UserUseCase{
		userRepo: repo,
	}
}

// UserSignup implements interfaces.UserUseCase.
func (c *UserUseCase) UserSignup(user helperStruct.UserReq) (response.UserData, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.UserData{}, err
	}
	user.Password = string(hash)
	userData, err := c.userRepo.UserSignUp(user)
	return userData, err
}

// UserLogin implements interfaces.UserUseCase.
func (c *UserUseCase) UserLogin(user helperStruct.LoginReq) (string, error) {
	var cfg config.Config
	userData, err := c.userRepo.UserLogin(user.Email)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password))
	fmt.Println(err)
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}
	if userData.IsBlocked {
		return "", fmt.Errorf("user is blocked")
	}
	claims := jwt.MapClaims{
		"id":  userData.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(cfg.SECRET))
	if err != nil {
		return "", err
	}

	return ss, nil
}
