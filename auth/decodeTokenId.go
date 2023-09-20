package auth

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

func DecodeTokenId(tokenString string) (float64, error) {

	jwtKey := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return 0, err
	}

	uidRaw, ok := claims["user_id"]
	if !ok {

	}

	uid, ok := uidRaw.(float64)
	if !ok {
		// Handle the case where the type assertion failed
		// This might occur if the "user_id" is not a number
	}

	return uid, nil
}
