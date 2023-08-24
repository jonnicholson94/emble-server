package utils

import (
	"errors"
	"os"

	"github.com/dgrijalva/jwt-go"
)

func DecodeTokenId(tokenString string) (string, error) {

	jwtKey := []byte(os.Getenv("JWT_KEY"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return "", err
	}

	var userID string
	switch v := claims["user_id"].(type) {
	case string:
		userID = v
	default:
		return "", errors.New("user_id is not a string")
	}

	return userID, nil
}
