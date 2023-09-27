package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func CreateToken(user_id int, first_name string, last_name string, auth_type string) (string, error) {

	jwtKey := []byte(os.Getenv("JWT_KEY"))

	claims := jwt.MapClaims{
		"user_id":    user_id,
		"first_name": first_name,
		"last_name":  last_name,
		"auth_type":  auth_type,
		"expiry":     time.Now().Add(time.Hour * 96).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
