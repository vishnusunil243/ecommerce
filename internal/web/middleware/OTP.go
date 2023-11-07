package middleware

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

var redisClient *redis.Client
var ctx = context.Background()

func init() {
	// Initialize the Redis client.
	redisClient = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("REDIS_ADDR"), // Update with your Redis server address.
		Password: "",                            // No password by default.
		DB:       0,                             // Use the default Redis database.
	})
}

// generateOTP generates a 6-digit random OTP.
func generateOTP() string {
	return strconv.Itoa(100000 + rand.Intn(900000))
}

func SendOTP(email string) error {

	// Create a new message.
	message := gomail.NewMessage()
	message.SetHeader("From", viper.GetString("SMTP_USER"))
	message.SetHeader("To", email)
	message.SetHeader("Subject", "OTP Verification")
	otp := generateOTP()

	// Set the message body
	message.SetBody("text/plain", "This will be expired in 5 minutes\nYour OTP is: "+otp)

	// Create an instance of the SMTP sender.
	dialer := gomail.NewDialer("smtp.gmail.com", 587, viper.GetString("SMTP_USER"), viper.GetString("SMTP_PASSWORD"))
	// Store the OTP in Redis with a short expiration time.
	otpKey := fmt.Sprintf("otp:%s", email)
	err := redisClient.Set(ctx, otpKey, otp, 300*time.Second).Err()
	if err != nil {
		log.Println("Failed to store OTP in Redis:", err)
		return err
	}
	// Send the email message.
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}
	return nil
}
func GetStoredOTP(email string) (string, error) {
	// Generate the Redis key used to store the OTP.
	otpKey := fmt.Sprintf("otp:%s", email)
	fmt.Println(otpKey)

	// Use the GET command to retrieve the OTP from Redis.
	otp, err := redisClient.Get(ctx, otpKey).Result()
	fmt.Println(otp)
	if err == redis.Nil {
		// OTP not found in Redis, return an error.
		return "", fmt.Errorf("OTP not found")
	} else if err != nil {
		// Handle other Redis errors.
		return "", err
	}

	return otp, nil
}
func VerifyOTP(email, otp string) bool {
	storedOTP, err := GetStoredOTP(email)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if otp == storedOTP {
		return true
	}
	return false
}
