package crud

import (
	"database/sql"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type FinalSurvey struct {
	ID           string           `json:"id"`
	Status       string           `json:"status"`
	PrototypeUrl string           `json:"prototype_url"`
	Questions    []SurveyQuestion `json:"questions"`
}

type JoinedSurvey struct {
	ID            string         `json:"id"`
	Status        string         `json:"status"`
	PrototypeUrl  string         `json:"prototype_url"`
	QuestionID    sql.NullString `json:"question_id"`
	QuestionTitle sql.NullString `json:"question_title"`
	QuestionType  sql.NullString `json:"question_type"`
	QuestionIndex sql.NullInt32  `json:"question_index"`
}

type SurveyQuestion struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Index int32  `json:"index"`
}

func FetchSurveyDetails(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	num, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Initialise the database

	db := utils.GetDB()

	query := "SELECT research.id, research.status, research.prototype_url, questions.id, questions.title, questions.type, questions.index FROM research LEFT JOIN questions ON research.id = questions.research_id WHERE research.id = $1"

	var finalSurvey FinalSurvey

	rows, err := db.Query(query, num)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for rows.Next() {
		var result JoinedSurvey

		scanErr := rows.Scan(
			&result.ID,
			&result.Status,
			&result.PrototypeUrl,
			&result.QuestionID,
			&result.QuestionTitle,
			&result.QuestionType,
			&result.QuestionIndex,
		)

		if scanErr != nil {
			fmt.Println(scanErr)
			http.Error(w, scanErr.Error(), http.StatusInternalServerError)
			return
		}

		question := SurveyQuestion{
			ID:    result.QuestionID.String,
			Title: result.QuestionTitle.String,
			Type:  result.QuestionType.String,
			Index: result.QuestionIndex.Int32,
		}

		if finalSurvey.ID != result.ID {
			finalSurvey.ID = result.ID
			finalSurvey.Status = result.Status
			finalSurvey.PrototypeUrl = result.PrototypeUrl
		}

		// Check if the question already exists in finalResearch
		questionExists := false
		for _, q := range finalSurvey.Questions {
			if q.ID == question.ID {
				questionExists = true
				break
			}
		}

		// If the question doesn't exist, append it to finalResearch
		if !questionExists {
			finalSurvey.Questions = append(finalSurvey.Questions, question)
		}
	}

	json, err := json.Marshal(finalSurvey)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "There's been a problem processing the json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
