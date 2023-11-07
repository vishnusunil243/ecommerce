package interfaces

import (
	"main.go/internal/domain"
)

type AdminRepository interface {
	AdminLogin(email string) (domain.Admins, error)
}
