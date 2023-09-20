package research

import (
	"database/sql"
	"emble-server/auth"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type FinalResearch struct {
	ID               string     `json:"ID"`
	Title            string     `json:"title"`
	Description      string     `json:"description"`
	Status           string     `json:"status"`
	Limit            int        `json:"limit"`
	PrototypeUrl     string     `json:"prototype_url"`
	UserId           string     `json:"user_id"`
	Intro            bool       `json:"intro"`
	IntroTitle       string     `json:"intro_title"`
	IntroDescription string     `json:"intro_description"`
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Questions        []Question `json:"questions"`
	Comments         []Comment  `json:"comments"`
}

type JoinedResearch struct {
	ID                 string         `json:"research_id"`
	Title              string         `json:"research_title"`
	Description        sql.NullString `json:"research_description"`
	Status             string         `json:"research_status"`
	Limit              int            `json:"research_limit"`
	PrototypeUrl       sql.NullString `json:"research_prototype_url"`
	UserId             string         `json:"research_user_id"`
	Intro              sql.NullBool   `json:"research_intro"`
	IntroTitle         sql.NullString `json:"research_intro_title"`
	IntroDescription   sql.NullString `json:"research_intro_description"`
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
	OptionID           sql.NullString `json:"option_id"`
	OptionContent      sql.NullString `json:"option_content"`
	OptionQuestionId   sql.NullString `json:"option_question_id"`
	OptionIndex        sql.NullInt32  `json:"option_index"`
	OptionResearchId   sql.NullString `json:"option_research_id"`
	FirstName          string         `json:"first_name"`
	LastName           string         `json:"last_name"`
}

type Question struct {
	QuestionID         string          `json:"question_id"`
	QuestionTitle      string          `json:"question_title"`
	QuestionType       string          `json:"question_type"`
	QuestionResearchId string          `json:"question_research_id"`
	QuestionIndex      int             `json:"question_index"`
	QuestionOptions    []FetchedOption `json:"question_options"`
}

type Comment struct {
	CommentID         string `json:"comment_id"`
	CommentContent    string `json:"comment_content"`
	CommentUserId     string `json:"comment_user_id"`
	CommentTimestamp  int    `json:"comment_timestamp"`
	CommentResearchId string `json:"comment_research_id"`
}

type FetchedOption struct {
	OptionID         string `json:"option_id"`
	OptionContent    string `json:"option_content"`
	OptionQuestionId string `json:"option_question_id"`
	OptionIndex      int    `json:"option_index"`
	OptionResearchId string `json:"option_research_id"`
}

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func FetchSingleResearch(w http.ResponseWriter, r *http.Request) {

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

	// Initialise the database

	db := utils.GetDB()

	// Define the query

	query := "SELECT research.*, questions.*, comments.*, options.*, users.first_name, users.last_name FROM research LEFT JOIN questions ON research.research_id = questions.question_research_id LEFT JOIN comments ON research.research_id = comments.comment_research_id LEFT JOIN options ON research.research_id = options.option_research_id LEFT JOIN users ON research.research_user_id = users.id WHERE research.research_id = $1;"

	var finalResearch FinalResearch

	rows, err := db.Query(query, id)

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
			&result.Intro,
			&result.IntroTitle,
			&result.IntroDescription,
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
			&result.OptionID,
			&result.OptionContent,
			&result.OptionQuestionId,
			&result.OptionIndex,
			&result.OptionResearchId,
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

		option := FetchedOption{
			OptionID:         result.OptionID.String,
			OptionContent:    result.OptionContent.String,
			OptionQuestionId: result.OptionQuestionId.String,
			OptionIndex:      int(result.OptionIndex.Int32),
			OptionResearchId: result.OptionResearchId.String,
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

		if finalResearch.ID != result.ID {
			finalResearch.ID = result.ID
			finalResearch.Title = result.Title
			finalResearch.Description = result.Description.String
			finalResearch.Limit = result.Limit
			finalResearch.Status = result.Status
			finalResearch.PrototypeUrl = result.PrototypeUrl.String
			finalResearch.UserId = result.UserId
			finalResearch.Intro = result.Intro.Bool
			finalResearch.IntroTitle = result.IntroTitle.String
			finalResearch.IntroDescription = result.IntroDescription.String
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

		for i, q := range finalResearch.Questions {
			if q.QuestionID == question.QuestionID && q.QuestionID == option.OptionQuestionId {
				// Check if the option already exists in QuestionOptions
				optionExists := false
				for _, o := range q.QuestionOptions {
					if o.OptionID == option.OptionID {
						optionExists = true
						break
					}
				}

				// If the option doesn't exist, append it to QuestionOptions
				if !optionExists {
					finalResearch.Questions[i].QuestionOptions = append(finalResearch.Questions[i].QuestionOptions, option)
				}
			}
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

	json, _ := json.Marshal(finalResearch)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
