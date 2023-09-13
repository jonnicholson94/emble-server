package auth

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type SignInUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	var user SignInUser

	err := json.NewDecoder(r.Body).Decode(&user)

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

	db := utils.GetDB()

	query := "SELECT first_name, last_name, id, email, password FROM users WHERE email = $1"

	var fetchedUser TokenUser

	queryError := db.QueryRow(query, user.Email).Scan(
		&fetchedUser.FirstName,
		&fetchedUser.LastName,
		&fetchedUser.ID,
		&fetchedUser.Email,
		&fetchedUser.Password,
	)

	if queryError != nil {
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

	hashErr := bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(user.Password))

	if hashErr != nil {
		customErr := CustomError{
			Message: "Invalid email or password",
			Status:  http.StatusBadRequest,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	fmt.Println("Successfully signed in.")

	// Generate JWT

	token, err := utils.CreateToken(fetchedUser.ID, fetchedUser.FirstName, fetchedUser.LastName)

	if err != nil {
		customErr := CustomError{
			Message: "Error creating token",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	// Return JWT to FE

	json, _ := json.Marshal(token)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
