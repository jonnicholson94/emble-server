package research

import (
	"database/sql"
	"emble-server/auth"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type FetchedResearch struct {
	ResearchID               string         `json:"research_id"`
	ResearchTitle            string         `json:"research_title"`
	ResearchDescription      sql.NullString `json:"research_description"`
	ResearchStatus           string         `json:"research_status"`
	ResearchLimit            string         `json:"research_limit"`
	ResearchPrototypeUrl     sql.NullString `json:"research_prototype_url"`
	ResearchUserId           int            `json:"research_user_id"`
	ResearchIntro            sql.NullBool   `json:"research_intro"`
	ResearchIntroTitle       sql.NullString `json:"research_intro_title"`
	ResearchIntroDescription sql.NullString `json:"research_intro_description"`
}

func FetchResearch(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	tokenErr := auth.ValidateToken(token)

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

	uid, err := auth.DecodeTokenId(token)

	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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
			&result.ResearchID,
			&result.ResearchTitle,
			&result.ResearchDescription,
			&result.ResearchStatus,
			&result.ResearchLimit,
			&result.ResearchPrototypeUrl,
			&result.ResearchUserId,
			&result.ResearchIntro,
			&result.ResearchIntroTitle,
			&result.ResearchIntroDescription,
		)

		if err != nil {
			fmt.Println(err)
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
		fmt.Println(err)
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
