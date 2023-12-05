package repository

import (
	"fmt"

	"gorm.io/gorm"
	"main.go/internal/repository/interfaces"
)

type ReferralDatabase struct {
	DB *gorm.DB
}

func NewReferralRepo(DB *gorm.DB) interfaces.ReferralRepository {
	return &ReferralDatabase{
		DB: DB,
	}
}

// AddReferral implements interfaces.ReferralRepository.
func (c *ReferralDatabase) AddReferral(userId int, referralId string) error {
	addReferral := `INSERT INTO referrals(user_id,referral_id) VALUES ($1,$2) `
	err := c.DB.Exec(addReferral, userId, referralId).Error
	if err != nil {
		return fmt.Errorf("error adding a referral id")
	}
	return nil
}

// ReferralOffer implements interfaces.ReferralRepository.
func (r *ReferralDatabase) ReferralOffer(referralId string, userId int) error {
	tx := r.DB.Begin()
	var exists bool
	tx.Raw(`SELECT EXISTS (SELECT 1 FROM user_referrals WHERE user_id=?)`, userId).Scan(&exists)
	if exists {
		tx.Rollback()
		return fmt.Errorf("you have already redeemed a referral offer please refer others to earn more exciting rewards")
	}
	var referredBy uint
	tx.Raw(`SELECT user_id FROM referrals WHERE referral_id=?`, referralId).Scan(&referredBy)
	if referredBy == 0 {
		tx.Rollback()
		return fmt.Errorf("invalid referralId Please enter a valid referral id")
	}
	if referredBy == uint(userId) {
		tx.Rollback()
		return fmt.Errorf("you can't refer yourselves")
	}
	err := tx.Exec(`INSERT INTO user_referrals (user_id,referred_by) VALUES ($1,$2)`, userId, referredBy).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var walletAmount int
	getWalletAmount := `SELECT amount FROM wallets WHERE user_id=$1`
	err = tx.Raw(getWalletAmount, userId).Scan(&walletAmount).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error retrieving amount from wallet")
	}
	insertWalletHistory := `INSERT INTO wallet_histories (recent_transaction,user_id,balance,time) VALUES ($1,$2,$3,NOW())`
	walletHistory := fmt.Sprintf("%d + 20", walletAmount)
	err = tx.Exec(insertWalletHistory, walletHistory, userId, (walletAmount + 20)).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting wallet history")
	}
	updateWallets := `UPDATE wallets SET amount=amount+20 WHERE user_id=?`
	err = tx.Exec(updateWallets, userId).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	var walletAmountOfReferrer int
	getWalletAmountOfReferrer := `SELECT amount FROM wallets WHERE user_id=$1`
	err = tx.Raw(getWalletAmountOfReferrer, referredBy).Scan(&walletAmountOfReferrer).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error retrieving amount from wallet")
	}
	insertWalletHistoryforReferrer := `INSERT INTO wallet_histories (recent_transaction,user_id,balance,time) VALUES ($1,$2,$3,NOW())`
	walletHistoryForReferrer := fmt.Sprintf("%d + 50", walletAmountOfReferrer)
	err = tx.Exec(insertWalletHistoryforReferrer, walletHistoryForReferrer, referredBy, (walletAmountOfReferrer + 50)).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting wallet history")
	}
	updateWalletsForReferrer := `UPDATE wallets SET amount=amount+50 WHERE user_id=?`
	err = tx.Exec(updateWalletsForReferrer, referredBy).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating wallet of referrer")
	}
	if err != nil {
		return err
	}
	if err = tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
