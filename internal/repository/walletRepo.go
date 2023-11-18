package repository

import (
	"gorm.io/gorm"
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
