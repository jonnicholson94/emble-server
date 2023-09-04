package crud

import (
	"database/sql"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type FinalResearch struct {
	ID           string     `json:"ID"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Status       string     `json:"status"`
	Limit        int        `json:"limit"`
	PrototypeUrl string     `json:"prototype_url"`
	UserId       string     `json:"user_id"`
	FirstName    string     `json:"first_name"`
	LastName     string     `json:"last_name"`
	Questions    []Question `json:"questions"`
	Comments     []Comment  `json:"comments"`
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
	CommentID          sql.NullString `json:"comment_id"`
	CommentContent     sql.NullString `json:"comment_content"`
	CommentUserId      sql.NullString `json:"comment_user_id"`
	CommentTimestamp   sql.NullInt64  `json:"comment_timestamp"`
	CommentResearchId  sql.NullString `json:"comment_research_id"`
	FirstName          string         `json:"first_name"`
	LastName           string         `json:"last_name"`
}

type Question struct {
	QuestionID         string `json:"question_id"`
	QuestionTitle      string `json:"question_title"`
	QuestionType       string `json:"question_type"`
	QuestionResearchId string `json:"question_research_id"`
	QuestionIndex      int    `json:"question_index"`
}

type Comment struct {
	CommentID         string `json:"comment_id"`
	CommentContent    string `json:"comment_content"`
	CommentUserId     string `json:"comment_user_id"`
	CommentTimestamp  int    `json:"comment_timestamp"`
	CommentResearchId string `json:"comment_research_id"`
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
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

	num, err := strconv.Atoi(id)

	// Initialise the database

	db := utils.GetDB()

	// Define the query

	query := "SELECT research.*, questions.*, comments.*, users.first_name, users.last_name FROM research LEFT JOIN questions ON research.id = questions.research_id LEFT JOIN comments ON research.id = comments.research_id LEFT JOIN users ON research.user_id = users.id WHERE research.id = $1;"

	var finalResearch FinalResearch

	rows, err := db.Query(query, num)

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
			&result.CommentID,
			&result.CommentContent,
			&result.CommentUserId,
			&result.CommentTimestamp,
			&result.CommentResearchId,
			&result.FirstName,
			&result.LastName,
		)

		// Append the question to the associated research's Questions slice
		question := Question{
			QuestionID:         result.QuestionID.String,
			QuestionTitle:      result.QuestionTitle.String,
			QuestionType:       result.QuestionType.String,
			QuestionResearchId: result.QuestionResearchId.String,
			QuestionIndex:      int(result.QuestionIndex.Int32),
		}

		comment := Comment{
			CommentID:         result.CommentID.String,
			CommentContent:    result.CommentContent.String,
			CommentUserId:     result.CommentUserId.String,
			CommentTimestamp:  int(result.CommentTimestamp.Int64),
			CommentResearchId: result.CommentResearchId.String,
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
			finalResearch.FirstName = result.FirstName
			finalResearch.LastName = result.LastName
		}

		// Check if the question already exists in finalResearch
		questionExists := false
		for _, q := range finalResearch.Questions {
			if q.QuestionID == question.QuestionID {
				questionExists = true
				break
			}
		}

		// If the question doesn't exist, append it to finalResearch
		if !questionExists {
			finalResearch.Questions = append(finalResearch.Questions, question)
		}

		// Check if the comment already exists in finalResearch
		commentExists := false
		for _, c := range finalResearch.Comments {
			if c.CommentID == comment.CommentID {
				commentExists = true
				break
			}
		}

		// If the comment doesn't exist, append it to finalResearch
		if !commentExists {
			finalResearch.Comments = append(finalResearch.Comments, comment)
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
