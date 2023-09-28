package auth

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Create the user struct
type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type ExistingUser struct {
	Email string `json:"email"`
}

type InsertedUser struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type CustomError struct {
	Message string `json:"message"`
	Status  int32  `json:"status"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Pull data from request body

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		customErr := CustomError{
			Message: "Error decoding request body",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	validPass := utils.ValidatePassword(user.Password)

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

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

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

	user.Password = string(hash)

	db := utils.GetDB()

	// Check if user already in database

	var existingUser ExistingUser

	existingQuery := "SELECT email FROM users WHERE email = $1"

	existingError := db.QueryRow(existingQuery, user.Email).Scan(&existingUser.Email)

	if existingError == nil {
		customErr := CustomError{
			Message: "User already exists",
			Status:  http.StatusBadRequest,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	// Save data in database

	query := "INSERT INTO users (first_name, last_name, email, password, premium) VALUES ($1, $2, $3, $4, false)"

	data, err := db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password)

	fmt.Println(data)

	var insertedUser InsertedUser

	selectQuery := "SELECT id, first_name, last_name FROM users WHERE email = $1"

	// Not getting scanned into inserted user struct

	err = db.QueryRow(selectQuery, user.Email).Scan(
		&insertedUser.ID,
		&insertedUser.FirstName,
		&insertedUser.LastName,
	)

	if err != nil {
		fmt.Println(err.Error())
		customErr := CustomError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	// Generate JWT

	token, err := utils.CreateToken(insertedUser.ID, insertedUser.FirstName, insertedUser.LastName, "password")

	if err != nil {
		customErr := CustomError{
			Message: "Failed to create user token",
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
