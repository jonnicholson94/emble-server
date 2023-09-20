package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(tokenString string) error {

	jwtKey := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		fmt.Println("Token isn't valid")
		return errors.New("Token isn't valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("Unable to map claims")
		return errors.New("Unable to map claims")
	}

	expiryClaim, ok := claims["expiry"].(float64)
	if !ok {
		fmt.Println("Invalid expiry date")
		return errors.New("Invalid expiry date")
	}

	expiryTime := time.Unix(int64(expiryClaim), 0)
	currentTime := time.Now()

	if currentTime.After(expiryTime) {
		fmt.Println("Auth token expired")
		return errors.New("Auth token expired")
	}

	return nil
}
