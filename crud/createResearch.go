package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type Research struct {
	ResearchId           string        `json:"research_id"`
	ResearchTitle        string        `json:"research_title"`
	ResearchDescription  string        `json:"research_description"`
	ResearchStatus       string        `json:"research_status"`
	ResearchLimit        int           `json:"research_limit"`
	ResearchPrototypeUrl string        `json:"research_prototype_url"`
	ResearchQuestions    []NewQuestion `json:"research_questions"`
}

func CreateResearch(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		fmt.Println(tokenErr.Error())
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
		fmt.Println(err.Error())
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

	query := "INSERT INTO research (research_id, research_title, research_description, research_status, research_limit, research_prototype_url, research_user_id) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	// Decode the token to get the user id

	uid, err := utils.DecodeTokenId(tk)

	if err != nil {
		fmt.Println(err.Error())
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

	_, queryErr := db.Exec(query, nr.ResearchId, nr.ResearchTitle, nr.ResearchDescription, nr.ResearchStatus, nr.ResearchLimit, nr.ResearchPrototypeUrl, uid)

	if queryErr != nil {
		fmt.Println(queryErr)
		customErr := CustomError{
			Message: "Query error",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	// fmt.Println(lastInsertID)

	if len(nr.ResearchQuestions) > 0 {
		questionQuery := "INSERT INTO questions (question_id, question_title, question_type, question_research_id, question_index) VALUES ($1, $2, $3, $4, $5)"

		for _, question := range nr.ResearchQuestions {
			_, err := db.Exec(questionQuery, question.QuestionId, question.QuestionTitle, question.QuestionType, nr.ResearchId, question.QuestionIndex)
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

			if len(question.QuestionOptions) > 0 {
				optionQuery := "INSERT INTO options (option_id, option_content, option_question_id, option_index, option_research_id) VALUES ($1, $2, $3, $4, $5)"
				for _, option := range question.QuestionOptions {
					_, err := db.Exec(optionQuery, option.OptionId, option.OptionContent, option.OptionQuestionID, option.OptionIndex, nr.ResearchId)
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
			}

		}
	}

	res, err := json.Marshal("Successfully created the research")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
