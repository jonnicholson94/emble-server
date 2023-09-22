package auth

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewUserDetails struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

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

	var newDetails NewUserDetails

	decodeErr := json.NewDecoder(r.Body).Decode(&newDetails)

	if decodeErr != nil {
		fmt.Println(err.Error())
		customErr := CustomError{
			Message: "Unable to decode body",
			Status:  http.StatusBadRequest,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	db := utils.GetDB()

	query := "UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3;"

	_, updateErr := db.Exec(query, newDetails.FirstName, newDetails.LastName, user_id)

	if updateErr != nil {
		fmt.Println(updateErr.Error())
		customErr := CustomError{
			Message: "Unable to update user details",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	json, _ := json.Marshal("Successfully updated users details")

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
