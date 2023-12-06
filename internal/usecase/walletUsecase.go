package usecase

import (
	"main.go/internal/common/helperStruct"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type walletUseCase struct {
	walletRepo interfaces.WalletRepository
}

func NewWalletUseCase(walletRepo interfaces.WalletRepository) services.WalletUseCase {
	return &walletUseCase{
		walletRepo: walletRepo,
	}
}

// CreateWallet implements interfaces.WalletUseCase.
func (w *walletUseCase) CreateWallet(userId int) error {
	err := w.walletRepo.CreateWallet(userId)
	return err
}

// DisplayWallet implements interfaces.WalletUseCase.
func (w *walletUseCase) DisplayWallet(userId int) (response.Wallet, error) {
	wallet, err := w.walletRepo.DisplayWallet(userId)
	return wallet, err
}

// WalletHistory implements interfaces.WalletUseCase.
func (w *walletUseCase) WalletHistory(userId int, queryParams helperStruct.QueryParams) ([]response.WalletHistories, int, error) {
	walletHistory, totalCount, err := w.walletRepo.WalletHistory(userId, queryParams)
	return walletHistory, totalCount, err
}
