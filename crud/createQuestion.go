package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewQuestion struct {
	QuestionId         string   `json:"question_id"`
	QuestionTitle      string   `json:"question_title"`
	QuestionType       string   `json:"question_type"`
	QuestionOptions    []Option `json:"question_options"`
	QuestionResearchId string   `json:"question_research_id"`
	QuestionIndex      int      `json:"question_index"`
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
		fmt.Println(err)
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

	insertQuery := "INSERT INTO questions (question_id, question_title, question_type, question_research_id, question_index) VALUES ($1, $2, $3, $4, $5)"

	data, err := db.Exec(insertQuery, nq.QuestionId, nq.QuestionTitle, nq.QuestionType, nq.QuestionResearchId, nq.QuestionIndex)

	fmt.Println(data)

	if err != nil {
		fmt.Println(err)
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
