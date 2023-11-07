package middleware

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"main.go/internal/infrastructure/config"
)

func ValidateToken(token string) (int, error) {
	var cfg config.Config
	Tokenvalue, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {

		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(cfg.SECRET), nil
	})

	if err != nil {
		return 0, err
	}

	if Tokenvalue == nil || !Tokenvalue.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	var parsedID interface{}
	if claims, ok := Tokenvalue.Claims.(jwt.MapClaims); ok && Tokenvalue.Valid {
		parsedID = claims["id"]
		//Check the expir
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return 0, fmt.Errorf("token expired")
		}
		// fmt.Println(claims["exp"])
	}
	value, ok := parsedID.(float64)
	if !ok {
		return 0, fmt.Errorf("expected an int value, but got %T", parsedID)
	}

	id := int(value)
	return id, err
}
