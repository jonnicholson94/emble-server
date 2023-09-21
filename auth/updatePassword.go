package auth

import (
	"database/sql"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type SubmittedPassword struct {
	Password string `json:"password"`
	ID       string `json:"id"`
}

type FetchedID struct {
	ResetToken string `json:"reset_token"`
}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {

	var password SubmittedPassword

	err := json.NewDecoder(r.Body).Decode(&password)

	if err != nil {
		customErr := CustomError{
			Message: "Failed to decode body",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	// Check the token is valid

	db := utils.GetDB()

	checkQuery := "SELECT users.reset_token FROM users WHERE reset_token = $1;"

	row := db.QueryRow(checkQuery, password.ID)

	var fetchedID FetchedID

	fetchErr := row.Scan(&fetchedID.ResetToken)
	if fetchErr != nil {
		if fetchErr == sql.ErrNoRows {
			// No rows were found
			fmt.Println(fetchErr.Error())

			customErr := CustomError{
				Message: "Invalid token",
				Status:  http.StatusBadRequest,
			}

			errJSON, _ := json.Marshal(customErr)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errJSON)
			return
		} else {
			// Handle other errors
			fmt.Println(fetchErr.Error())

			customErr := CustomError{
				Message: "There was a problem while fetching the data",
				Status:  http.StatusInternalServerError,
			}

			errJSON, _ := json.Marshal(customErr)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errJSON)
			return
		}
	}

	// Update user's password and reset token

	validPass := utils.ValidatePassword(password.Password)

	if !validPass {
		customErr := CustomError{
			Message: "Invalid password provided",
			Status:  http.StatusBadRequest,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	// Hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(password.Password), bcrypt.DefaultCost)

	if err != nil {
		customErr := CustomError{
			Message: "Failed to hash password",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	password.Password = string(hash)

	updateQuery := "UPDATE users SET password = $1, reset_token = $2 WHERE reset_token = $3;"

	_, insertErr := db.Exec(updateQuery, password.Password, "", password.ID)

	if insertErr != nil {
		fmt.Println(insertErr.Error())
		customErr := CustomError{
			Message: "Error inserting to database",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	// Return success

	json, _ := json.Marshal("Successfully updated users password")

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
