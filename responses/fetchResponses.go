package responses

import (
	"emble-server/auth"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type CustomError struct {
	Message string `json:"message"`
	Status  int32  `json:"status"`
}

type JoinedResponse struct {
	ResearchTitle      string `json:"research_title"`
	QuestionID         string `json:"question_id"`
	QuestionTitle      string `json:"question_title"`
	QuestionType       string `json:"question_type"`
	ResponseID         string `json:"response_id"`
	ResponseQuestionID string `json:"response_question_id"`
	ResponseAnswer     string `json:"response_answer"`
}

type Response struct {
	ResponseID         string `json:"response_id"`
	ResponseQuestionID string `json:"response_question_id"`
	ResponseAnswer     string `json:"response_answer"`
}

type Question struct {
	QuestionID    string     `json:"question_id"`
	QuestionTitle string     `json:"question_title"`
	QuestionType  string     `json:"question_type"`
	Responses     []Response `json:"question_responses"`
}

type FinalResponse struct {
	ResearchTitle string     `json:"research_title"`
	Questions     []Question `json:"research_questions"`
}

func FetchResponses(w http.ResponseWriter, r *http.Request) {

	// Check token and fetch IDtk := r.Header.Get("Authorization")

	tk := r.Header.Get("Authorization")

	tokenErr := auth.ValidateToken(tk)

	if tokenErr != nil {
		fmt.Println(tokenErr)
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

	db := utils.GetDB()

	selectQuery := "SELECT research.research_title, questions.question_id, questions.question_title, questions.question_type, responses.response_id, responses.response_question_id, responses.response_answer FROM research LEFT JOIN questions ON research.research_id = questions.question_research_id LEFT JOIN responses ON research.research_id = responses.response_research_id WHERE research.research_id = $1;"

	// Conduct query

	var finalResponse FinalResponse

	rows, err := db.Query(selectQuery, id)

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

	// Scan rows

	for rows.Next() {

		var result JoinedResponse

		scanErr := rows.Scan(
			&result.ResearchTitle,
			&result.QuestionID,
			&result.QuestionTitle,
			&result.QuestionType,
			&result.ResponseID,
			&result.ResponseQuestionID,
			&result.ResponseAnswer,
		)

		question := Question{
			QuestionID:    result.QuestionID,
			QuestionTitle: result.QuestionTitle,
			QuestionType:  result.QuestionType,
		}

		response := Response{
			ResponseID:         result.ResponseID,
			ResponseQuestionID: result.ResponseQuestionID,
			ResponseAnswer:     result.ResponseAnswer,
		}

		if scanErr != nil {
			fmt.Println(scanErr)
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

		// Manipulate data

		if finalResponse.ResearchTitle != result.ResearchTitle {
			finalResponse.ResearchTitle = result.ResearchTitle
		}

		questionExists := false

		for _, q := range finalResponse.Questions {
			if q.QuestionID == question.QuestionID {
				questionExists = true
				break
			}
		}

		if !questionExists {
			finalResponse.Questions = append(finalResponse.Questions, question)
		}

		for i, q := range finalResponse.Questions {
			if q.QuestionID == response.ResponseQuestionID {
				responseExists := false

				for _, r := range finalResponse.Questions[i].Responses {
					if r.ResponseID == response.ResponseID {
						responseExists = true
						break
					}
				}

				if !responseExists {
					finalResponse.Questions[i].Responses = append(finalResponse.Questions[i].Responses, response)
				}

				break
			}
		}

	}

	// Send back to FE

	json, _ := json.Marshal(finalResponse)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
