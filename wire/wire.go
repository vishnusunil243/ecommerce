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
		repository.NewWalletRepo,
		repository.NewUserRepo,
		repository.NewAdminRepo,
		repository.NewProductRepo,
		repository.NewSuperRepo,
		repository.NewCartRepo,
		repository.NewOrderRepo,
		usecase.NewUserUsecase,
		usecase.NewAdminUsecase,
		usecase.NewProductUsecase,
		usecase.NewSuperAdminUsecase,
		usecase.NewCartUseCase,
		usecase.NewOrderUseCase,
		usecase.NewWalletUseCase,
		handler.NewUserHandler,
		handler.NewAdminHandler,
		handler.NewProductHandler,
		handler.NewSuperAdminHandler,
		handler.NewCartHandler,
		handler.NewOrderHandler,
		http.NewServerHTTP,
	)
	return &http.ServerHTTP{}, nil
}
