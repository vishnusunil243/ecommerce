package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/response"
	"main.go/internal/domain"
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
	// Fetch the maximum existing id from the table
	var maxID int
	err := w.DB.Model(&domain.Wallet{}).Select("COALESCE(MAX(id), 0)").Row().Scan(&maxID)
	if err != nil {
		return err
	}

	// Increment the maximum id for the new wallet
	newID := maxID + 1
	createWallet := `INSERT INTO wallets(id,user_id,amount) VALUES ($1,$2,0)`
	err = w.DB.Exec(createWallet, newID, userId).Error
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
