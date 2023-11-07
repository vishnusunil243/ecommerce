package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"main.go/internal/domain"
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
		&domain.Product{},
		&domain.Orders{},
	)
	return db, err
}
