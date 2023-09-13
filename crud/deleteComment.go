package crud

import (
	"emble-server/utils"
	"encoding/json"
	"net/http"
)

func DeleteComment(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
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

	db := utils.GetDB()

	id := r.URL.Query().Get("id")

	_, err := db.Exec("DELETE FROM comments WHERE id = $1", id)

	if err != nil {
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

	res, jsonErr := json.Marshal("Successfully deleted the comment")

	if jsonErr != nil {
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
