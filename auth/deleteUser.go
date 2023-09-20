package auth

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserToDelete struct {
	Email string `json:"email"`
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	var user UserToDelete

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

	query := "DELETE FROM users WHERE email = $1;"

	_, queryErr := db.Exec(query, user.Email)

	if queryErr != nil {
		fmt.Println(queryErr)
		customErr := CustomError{
			Message: "Failed to process request",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	res, _ := json.Marshal("Successfully deleted the user")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
