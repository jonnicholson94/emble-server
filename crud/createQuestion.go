package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewQuestion struct {
	Title      string `json:"title"`
	Type       string `json:"type"`
	ResearchId string `json:"research_id"`
	Index      int    `json:"index"`
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {

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

	var nq NewQuestion

	err := json.NewDecoder(r.Body).Decode(&nq)

	if err != nil {
		customErr := CustomError{
			Message: "Failed to decode data",
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

	insertQuery := "INSERT INTO questions (title, type, research_id, index) VALUES ($1, $2, $3, $4)"

	data, err := db.Exec(insertQuery, nq.Title, nq.Type, nq.ResearchId, nq.Index)

	fmt.Println(data)

	if err != nil {
		customErr := CustomError{
			Message: "Failed to insert data",
			Status:  http.StatusBadRequest,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	res, err := json.Marshal("Successfully created the question")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
