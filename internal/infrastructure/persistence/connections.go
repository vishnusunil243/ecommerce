package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"main.go/internal/domain"
	"main.go/internal/infrastructure/concurrency"
	"main.go/internal/infrastructure/config"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	config.LoadConfig()
	database_key := cfg.DB_KEY
	db, err := gorm.Open(postgres.Open(database_key), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("error connecting to database")
	}
	db.AutoMigrate(
		&domain.Admins{},
		&domain.Users{},
		&domain.UserInfo{},
		&domain.Address{},
		&domain.Product{},
		&domain.ProductItem{},
		&domain.PaymentType{},
		&domain.Category{},
		&domain.PaymentDetails{},
		&domain.PaymentStatus{},
		domain.Brand{},
		&domain.Orders{},
		&domain.OrderItem{},
		&domain.SuperAdmin{},
		&domain.AdminInfo{},
		&domain.ReportInfo{},
		&domain.Images{},
		&domain.Image_items{},
		&domain.Carts{},
		&domain.Wallet{},
		&domain.CartItem{},
		&domain.Coupon{},
		&domain.UserCoupons{},
		&domain.WalletHistories{},
		&domain.Wishlist{},
		&domain.Discount{},
		&domain.Referrals{},
		domain.UserReferrals{},
	)
	unblockUser := concurrency.NewConcurrency(db)

	// Start the UserStatusChecker goroutine
	unblockUser.Concurrency()
	return db, err
}
