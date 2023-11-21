package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/response"
	"main.go/internal/repository/interfaces"
)

type walletRepository struct {
	DB *gorm.DB
}

func NewWalletRepo(DB *gorm.DB) interfaces.WalletRepository {
	return &walletRepository{
		DB: DB,
	}
}

// CreateWallet implements interfaces.WalletRepository.
func (w *walletRepository) CreateWallet(userId int) error {
	createWallet := `INSERT INTO wallets(user_id,amount) VALUES ($1,0)`
	err := w.DB.Exec(createWallet, userId).Error
	return err
}

// DisplayWallet implements interfaces.WalletRepository.
func (w *walletRepository) DisplayWallet(userId int) (response.Wallet, error) {
	var wallet response.Wallet
	displayWalletForUser := `SELECT * FROM wallets WHERE user_id=$1`
	err := w.DB.Raw(displayWalletForUser, userId).Scan(&wallet).Error
	if err != nil {
		return response.Wallet{}, fmt.Errorf("error displaying wallet please contact the customer care")
	}

	return wallet, nil
}
