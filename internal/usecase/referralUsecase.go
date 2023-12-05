package usecase

import (
	"fmt"

	"main.go/internal/repository/interfaces"
	services "main.go/internal/usecase/interface"
)

type ReferralUseCase struct {
	referralRepo interfaces.ReferralRepository
}

func NewReferralUsecase(referralRepo interfaces.ReferralRepository) services.ReferralUseCase {
	return &ReferralUseCase{
		referralRepo: referralRepo,
	}
}

// AddReferral implements interfaces.ReferralUseCase.
func (r *ReferralUseCase) AddReferral(userId int, userMobile string) error {
	referralCode := fmt.Sprintf("pc%s/%d4u", userMobile, userId)
	err := r.referralRepo.AddReferral(userId, string(referralCode))
	return err
}

// ReferralOffer implements interfaces.ReferralUseCase.
func (r *ReferralUseCase) ReferralOffer(userId int, referralId string) error {
	err := r.referralRepo.ReferralOffer(referralId, userId)
	return err
}
