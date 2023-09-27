package oauth

import (
	"context"
	"database/sql"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type CallbackProperties struct {
	AuthUser string `json:"authuser"`
	Code     string `json:"code"`
	Prompt   string `json:"prompt"`
	Scope    string `json:"scope"`
}

type UserInfo struct {
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
}

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type CustomError struct {
	Message string `json:"message"`
	Status  int32  `json:"status"`
}

func ExchangeCode(w http.ResponseWriter, r *http.Request) {

	var properties CallbackProperties

	err := json.NewDecoder(r.Body).Decode(&properties)

	if err != nil {
		fmt.Println(err.Error())
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("BASE_URL") + "/auth/callback",
		Scopes: []string{
			"email",
			"profile",
			"offline_access",
		},
		Endpoint: google.Endpoint,
	}

	var context = context.TODO()

	token, err := config.Exchange(context, properties.Code)

	if err != nil {
		fmt.Println(err.Error())
	}

	// Get the user's details

	res, err := GetUserInfo(token.AccessToken)

	if err != nil {
		fmt.Println("Get user info error")
		fmt.Println(err)
	}

	var user UserInfo

	decodeErr := json.NewDecoder(res.Body).Decode(&user)

	if decodeErr != nil {
		fmt.Println("Decode error")
		fmt.Println(decodeErr.Error())
	}

	// Check if user exists in database

	var existingUser User

	db := utils.GetDB()

	query := "SELECT id, first_name, last_name, email FROM users WHERE email = $1;"

	queryError := db.QueryRow(query, user.Email).Scan(
		&existingUser.ID,
		&existingUser.FirstName,
		&existingUser.LastName,
		&existingUser.Email,
	)

	if queryError != nil {
		fmt.Println(queryError.Error())
		// 	// Handle user not existing
		// 	// Create new user with appropriate details

		if queryError == sql.ErrNoRows {

			insertQuery := "INSERT INTO users (first_name, last_name, email, premium, google_access_token, google_refresh_token) VALUES ($1, $2, $3, $4, $5, $6)"

			_, err := db.Exec(insertQuery, user.GivenName, user.FamilyName, user.Email, false, token.AccessToken, token.RefreshToken)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			selectQuery := "SELECT id FROM users WHERE email = $1"

			selectErr := db.QueryRow(selectQuery, user.Email).Scan(&existingUser.ID)

			if selectErr != nil {
				fmt.Println(selectErr.Error())
				return
			}

		}

	} else {

		updateQuery := "UPDATE users SET google_access_token = $1, google_refresh_token = $2 WHERE id = $3;"

		_, err := db.Exec(updateQuery, token.AccessToken, token.RefreshToken, existingUser.ID)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

	}

	jwt, _ := utils.CreateToken(existingUser.ID, existingUser.FirstName, existingUser.LastName, "google")

	json, _ := json.Marshal(jwt)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
