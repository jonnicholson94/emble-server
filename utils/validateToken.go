package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(w http.ResponseWriter, r *http.Request) {
	jwtKey := []byte(os.Getenv("JWT_KEY"))
	tokenString := r.Header.Get("Authorization")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusUnauthorized)
		return
	}

	fmt.Println(claims["user_id"])
	fmt.Println(claims["first_name"])
	fmt.Println(claims["last_name"])
	fmt.Println(time.Now().Unix())
	fmt.Println(int(claims["expiry"].(float64)))

}
