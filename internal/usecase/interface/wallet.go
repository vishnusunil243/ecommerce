package interfaces

type WalletUseCase interface {
	CreateWallet(userId int) error
}
