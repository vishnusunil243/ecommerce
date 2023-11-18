package interfaces

type WalletRepository interface {
	CreateWallet(userId int) error
}
