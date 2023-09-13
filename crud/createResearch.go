package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type Research struct {
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Status       string        `json:"status"`
	Limit        int           `json:"limit"`
	PrototypeUrl string        `json:"prototype_url"`
	Questions    []NewQuestion `json:"questions"`
}

func CreateResearch(w http.ResponseWriter, r *http.Request) {

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

	var nr Research

	err := json.NewDecoder(r.Body).Decode(&nr)

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

	fmt.Println(nr)

	db := utils.GetDB()

	query := "INSERT INTO research (title, description, status, \"limit\", prototype_url, user_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"

	// Decode the token to get the user id

	uid, err := utils.DecodeTokenId(tk)

	if err != nil {
		customErr := CustomError{
			Message: "Failed to process request",
			Status:  http.StatusBadRequest,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	var lastInsertID int

	queryErr := db.QueryRow(query, nr.Title, nr.Description, nr.Status, nr.Limit, nr.PrototypeUrl, uid).Scan(&lastInsertID)

	if queryErr != nil {
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

	fmt.Println(lastInsertID)

	if len(nr.Questions) > 0 {
		questionQuery := "INSERT INTO questions (title, type, research_id, \"index\") VALUES ($1, $2, $3, $4)"

		for _, question := range nr.Questions {
			_, err := db.Exec(questionQuery, question.Title, question.Type, lastInsertID, question.Index)
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
		}
	}

	res, err := json.Marshal("Successfully created the research")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
