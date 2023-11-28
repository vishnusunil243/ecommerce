package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DB_KEY         string `mapstructure:"DB_KEY"`
	SMTP_SERVER    string `mapstructure:"SMTP_SERVER"`
	SMTP_PORT      string `mapstructure:"SMTP_PORT"`
	SMTP_PASSORD   string `mapstructure:"SMTP_PASSWORD"`
	SMTP_USER      string `mapstructure:"SMTP_USER"`
	REDIS_ADDR     string `mapstructure:"REDIS_ADDR"`
	SECRET         string `mapstructure:"SECRET"`
	RAZORPAYID     string `mapstructure:"RAZORPAY_ID"`
	RAZORPAYSECRET string `mapstructure:"RAZORPAY_SECRET"`
}

var envs = []string{
	"DB_KEY",
	"SMTP_SERVER",
	"SMTP_USER",
	"SMTP_PORT",
	"SMTP_PASSWORD",
	"REDIS_ADDR",
	"SECRET",
	"RAZORPAY_ID",
	"RAZORPAY_SECRET",
}

func LoadConfig() (Config, error) {
	var config Config

	viper.AddConfigPath("../../.env")
	viper.SetConfigFile("../../.env")
	viper.ReadInConfig()

	for _, env := range envs {
		if err := viper.BindEnv(env); err != nil {
			return config, err
		}
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	if err := validator.New().Struct(&config); err != nil {
		return config, err
	}

	return config, nil
}
