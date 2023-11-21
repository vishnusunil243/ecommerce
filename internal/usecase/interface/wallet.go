package interfaces

import "main.go/internal/common/response"

type WalletUseCase interface {
	CreateWallet(userId int) error
	DisplayWallet(userId int) (response.Wallet, error)
}
