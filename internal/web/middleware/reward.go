package middleware

import (
	"fmt"

	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendCouponEmail(userEmail, couponCode string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "your-email@example.com")
	m.SetHeader("To", userEmail)
	m.SetHeader("Subject", "Reward For Your Loyalty")
	m.SetBody("text/html", fmt.Sprintf("Dear user,<br>Thank you for choosing us everytime.We have seen that you have completed 10 orders with us as a recognition of loyalty please accept this coupon with code: <strong>%s</strong>", couponCode))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, viper.GetString("SMTP_USER"), viper.GetString("SMTP_PASSWORD"))

	// Send the email
	if err := dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
