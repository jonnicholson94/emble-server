package crud

import (
	"database/sql"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type FinalResearch struct {
	ID           string     `json:"ID"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Limit        int        `json:"limit"`
	PrototypeUrl string     `json:"prototype_url"`
	UserId       string     `json:"user_id"`
	Questions    []Question `json:"questions"`
}

type JoinedResearch struct {
	ID                 string         `json:"ID"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	Status             string         `json:"status"`
	Limit              int            `json:"limit"`
	PrototypeUrl       string         `json:"prototype_url"`
	UserId             string         `json:"user_id"`
	QuestionID         sql.NullString `json:"question_id"`
	QuestionTitle      sql.NullString `json:"question_title"`
	QuestionType       sql.NullString `json:"question_type"`
	QuestionResearchId sql.NullString `json:"question_research_id"`
	QuestionIndex      sql.NullInt32  `json:"question_index"`
}

type Question struct {
	QuestionID         string `json:"question_id"`
	QuestionTitle      string `json:"question_title"`
	QuestionType       string `json:"question_type"`
	QuestionResearchId string `json:"question_research_id"`
	QuestionIndex      int    `json:"question_index"`
}

func FetchSingleResearch(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		fmt.Println(tokenErr)
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")

	// Initialise the database

	db := utils.GetDB()

	// Define the query

	query := "SELECT research.*, questions.* FROM research LEFT JOIN questions ON research.id = questions.research_id WHERE research.id = $1;"

	var finalResearch FinalResearch

	rows, err := db.Query(query, id)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for rows.Next() {
		var result JoinedResearch

		scanErr := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Description,
			&result.Status,
			&result.Limit,
			&result.PrototypeUrl,
			&result.UserId,
			&result.QuestionID,
			&result.QuestionTitle,
			&result.QuestionType,
			&result.QuestionResearchId,
			&result.QuestionIndex,
		)

		// Append the question to the associated research's Questions slice
		question := Question{
			QuestionID:         result.QuestionID.String,
			QuestionTitle:      result.QuestionTitle.String,
			QuestionType:       result.QuestionType.String,
			QuestionResearchId: result.QuestionResearchId.String,
			QuestionIndex:      int(result.QuestionIndex.Int32),
		}

		if scanErr != nil {
			fmt.Println(scanErr)
			http.Error(w, scanErr.Error(), http.StatusInternalServerError)
			return
		}

		if finalResearch.ID != result.ID {
			finalResearch.ID = result.ID
			finalResearch.Title = result.Title
			finalResearch.Description = result.Description
			finalResearch.Limit = result.Limit
			finalResearch.Status = result.Status
			finalResearch.PrototypeUrl = result.PrototypeUrl
			finalResearch.UserId = result.UserId
			finalResearch.Questions = append(finalResearch.Questions, question)
		} else {
			finalResearch.Questions = append(finalResearch.Questions, question)
		}

	}

	json, err := json.Marshal(finalResearch)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "There's been a problem processing the json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
