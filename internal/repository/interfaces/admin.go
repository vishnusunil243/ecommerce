package interfaces

import (
	"main.go/internal/common/response"
	"main.go/internal/domain"
)

type AdminRepository interface {
	AdminLogin(email string) (domain.Admins, error)
	ListAllUsers() ([]response.UserDetails, error)
	DispalyUser(id int) (response.UserDetails, error)
}
