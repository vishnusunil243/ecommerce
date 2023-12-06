package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/common/helperStruct"
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

// WalletHistory implements interfaces.WalletRepository.
func (w *walletRepository) WalletHistory(userid int, queryParams helperStruct.QueryParams) ([]response.WalletHistories, int, error) {
	var walletHistories []response.WalletHistories
	walletHistory := `SELECT * FROM wallet_histories WHERE user_id=?`
	var count int
	getTotalCount := fmt.Sprintf("SELECT COUNT(*) FROM (%s%d)", walletHistory[:len(walletHistory)-1], userid)
	err := w.DB.Raw(getTotalCount).Scan(&count).Error
	if err != nil {
		return []response.WalletHistories{}, 0, err
	}
	walletHistory = fmt.Sprintf("%s ORDER BY time DESC", walletHistory)
	if queryParams.Limit != 0 && queryParams.Page != 0 {
		walletHistory = fmt.Sprintf("%s LIMIT %d OFFSET %d", walletHistory, queryParams.Limit, (queryParams.Page-1)*queryParams.Limit)
	}
	if queryParams.Limit == 0 || queryParams.Page == 0 {
		walletHistory = fmt.Sprintf("%s LIMIT 10 OFFSET 0", walletHistory)
	}
	err = w.DB.Raw(walletHistory, userid).Scan(&walletHistories).Error
	return walletHistories, count, err
}
