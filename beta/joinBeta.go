package beta

import (
	"emble-server/utils"
	"encoding/json"
	"net/http"
	"time"
)

type CustomError struct {
	Message string `json:"message"`
	Status  int32  `json:"status"`
}

type Email struct {
	Email string `json:"email"`
}

func JoinBeta(w http.ResponseWriter, r *http.Request) {

	var newEmail Email
	currentTime := time.Now()

	decodeErr := json.NewDecoder(r.Body).Decode(&newEmail)

	if decodeErr != nil {
		customErr := CustomError{
			Message: "Error decoding body",
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

	insertQuery := "INSERT INTO beta (email, timestamp) VALUES ($1, $2)"

	_, err := db.Exec(insertQuery, newEmail.Email, currentTime)

	if err != nil {
		customErr := CustomError{
			Message: "There's been a problem signing you up. You must have already joined.",
			Status:  http.StatusBadRequest,
		}

		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSON)
		return
	}

	res, _ := json.Marshal("Successfully signed up")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
