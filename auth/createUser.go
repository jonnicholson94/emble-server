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

func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Pull data from request body

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Incorrect details provided", http.StatusBadRequest)
		return
	}

	validPass := utils.ValidatePassword(user.Password)

	if !validPass {
		http.Error(w, "Invalid password provided. Make sure you enter between 6 and 20 characters, and the password contains at least one number and special character.", http.StatusBadRequest)
		return
	}

	// Hash the password

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user.Password = string(hash)

	db := utils.GetDB()

	// Check if user already in database

	var existingUser ExistingUser

	existingQuery := "SELECT email FROM users WHERE email = $1"

	existingError := db.QueryRow(existingQuery, user.Email).Scan(&existingUser.Email)

	if existingError == nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
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
	}

	// Generate JWT

	token, err := utils.CreateToken(insertedUser.ID, insertedUser.FirstName, insertedUser.LastName)

	if err != nil {
		http.Error(w, "There's been an issue generating your token", http.StatusInternalServerError)
		return
	}

	// Return JWT to FE

	json, err := json.Marshal(token)

	if err != nil {
		http.Error(w, "There's been an issue sending your token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
