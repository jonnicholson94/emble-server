package auth

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := ValidateToken(tk)

	if tokenErr != nil {
		fmt.Println(tokenErr.Error())
		customErr := CustomError{
			Message: "Invalid token",
			Status:  http.StatusUnauthorized,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJSON)
		return
	}

	user_id, err := DecodeTokenId(tk)

	if err != nil {
		fmt.Println(err.Error())
		customErr := CustomError{
			Message: "Invalid token",
			Status:  http.StatusUnauthorized,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJSON)
		return
	}

	var user UserDetails

	db := utils.GetDB()

	query := "SELECT first_name, last_name, email FROM users WHERE id = $1;"

	row := db.QueryRow(query, user_id)

	scanErr := row.Scan(
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)

	if scanErr != nil {
		customErr := CustomError{
			Message: "No user with this email exists",
			Status:  http.StatusBadRequest,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	json, _ := json.Marshal(user)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
