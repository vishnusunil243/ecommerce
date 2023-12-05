package interfaces

type ReferralUseCase interface {
	AddReferral(userId int, userMobile string) error
	ReferralOffer(userId int, referralId string) error
}
