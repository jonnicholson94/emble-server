package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	ResponseId           string `json:"answer_id"`
	ResponseAnswer       string `json:"answer_answer"`
	ResponseQuestionId   string `json:"answer_question_id"`
	ResponseQuestionType string `json:"answer_question_type"`
	ResponseResearchId   string `json:"answer_research_id"`
}

func CreateResponse(w http.ResponseWriter, r *http.Request) {

	var responses []Response

	err := json.NewDecoder(r.Body).Decode(&responses)

	if err != nil {
		customErr := CustomError{
			Message: "Failed to decode JSON",
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

	createQuery := "INSERT INTO responses (response_id, response_answer, response_question_id, response_question_type, response_research_id) VALUES ($1, $2, $3, $4, $5)"

	for _, response := range responses {
		_, err := db.Exec(createQuery, response.ResponseId, response.ResponseAnswer, response.ResponseQuestionId, response.ResponseQuestionType, response.ResponseResearchId)

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
	}

	res, _ := json.Marshal("Successfully created the response")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
