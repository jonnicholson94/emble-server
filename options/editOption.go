package options

import (
	"emble-server/auth"
	"emble-server/utils"
	"encoding/json"
	"net/http"
)

type UpdatedOption struct {
	OptionContent string `json:"option_content"`
}

func EditOption(w http.ResponseWriter, r *http.Request) {
	tk := r.Header.Get("Authorization")

	tokenErr := auth.ValidateToken(tk)

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

	id := r.URL.Query().Get("id")

	var option UpdatedOption

	err := json.NewDecoder(r.Body).Decode(&option)

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

	db := utils.GetDB()

	query := "UPDATE options SET option_content = $1 WHERE option_id = $2"

	_, dbErr := db.Exec(query, option.OptionContent, id)

	if dbErr != nil {
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

	res, err := json.Marshal("Successfully updated the comment")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
