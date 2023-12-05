package interfaces

type ReferralRepository interface {
	AddReferral(userId int, referralId string) error
	ReferralOffer(referralId string, userId int) error
}
