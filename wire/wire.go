package wire

import (
	"github.com/google/wire"
	"main.go/internal/infrastructure/config"
	db "main.go/internal/infrastructure/persistence"
	"main.go/internal/repository"
	"main.go/internal/usecase"
	http "main.go/internal/web"
	"main.go/internal/web/handler"
)

func InitializeAPI1(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(
		db.ConnectDatabase,
		repository.NewUserRepo,
		repository.NewAdminRepo,
		usecase.NewUserUsecase,
		usecase.NewAdminUsecase,
		handler.NewUserHandler,
		handler.NewAdminHandler,
		http.NewServerHTTP,
	)
	return &http.ServerHTTP{}, nil
}
