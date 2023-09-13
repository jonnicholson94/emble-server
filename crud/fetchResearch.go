package crud

import (
	"emble-server/utils"
	"encoding/json"
	"net/http"
)

type FetchedResearch struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Limit        string `json:"limit"`
	PrototypeUrl string `json:"prototype_url"`
	UserId       int    `json:"user_id"`
}

func FetchResearch(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(token)

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

	uid, err := utils.DecodeTokenId(token)

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

	query := "SELECT * FROM research WHERE research_user_id = $1"

	rows, err := db.Query(query, uid)

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

	var results []FetchedResearch

	for rows.Next() {
		var result FetchedResearch

		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Description,
			&result.Status,
			&result.Limit,
			&result.PrototypeUrl,
			&result.UserId,
		)

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

		results = append(results, result)
	}

	data, err := json.Marshal(results)

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

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
