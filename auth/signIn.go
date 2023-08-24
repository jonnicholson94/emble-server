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
		http.Error(w, "There was an error decoding the request body", http.StatusBadRequest)
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
		http.Error(w, "No account with this email exists", http.StatusBadRequest)
		return
	}

	hashErr := bcrypt.CompareHashAndPassword([]byte(fetchedUser.Password), []byte(user.Password))

	if hashErr != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	fmt.Println("Successfully signed in.")

	// Generate JWT

	token, err := utils.CreateToken(fetchedUser.ID, fetchedUser.FirstName, fetchedUser.LastName)

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
