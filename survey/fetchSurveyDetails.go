package survey

import (
	"database/sql"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type FinalSurvey struct {
	ResearchID               string           `json:"id"`
	ResearchStatus           string           `json:"status"`
	ResearchPrototypeUrl     string           `json:"prototype_url"`
	ResearchIntro            bool             `json:"intro"`
	ResearchIntroTitle       string           `json:"intro_title"`
	ResearchIntroDescription string           `json:"intro_description"`
	ResearchQuestions        []SurveyQuestion `json:"questions"`
}

type JoinedSurvey struct {
	ResearchID               string         `json:"id"`
	ResearchStatus           string         `json:"status"`
	ResearchPrototypeUrl     string         `json:"prototype_url"`
	ResearchIntro            bool           `json:"intro"`
	ResearchIntroTitle       string         `json:"intro_title"`
	ResearchIntroDescription string         `json:"intro_description"`
	QuestionID               sql.NullString `json:"question_id"`
	QuestionTitle            sql.NullString `json:"question_title"`
	QuestionType             sql.NullString `json:"question_type"`
	QuestionIndex            sql.NullInt32  `json:"question_index"`
	OptionID                 sql.NullString `json:"option_id"`
	OptionContent            sql.NullString `json:"option_content"`
	OptionQuestionId         sql.NullString `json:"option_question_id"`
	OptionIndex              sql.NullInt32  `json:"option_index"`
	OptionResearchId         sql.NullString `json:"option_research_id"`
}

type SurveyQuestion struct {
	QuestionID      string          `json:"question_id"`
	QuestionTitle   string          `json:"question_title"`
	QuestionType    string          `json:"question_type"`
	QuestionIndex   int32           `json:"question_index"`
	QuestionOptions []FetchedOption `json:"question_options"`
}

type FetchedOption struct {
	OptionID         string `json:"option_id"`
	OptionContent    string `json:"option_content"`
	OptionQuestionId string `json:"option_question_id"`
	OptionIndex      int    `json:"option_index"`
	OptionResearchId string `json:"option_research_id"`
}

func FetchSurveyDetails(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	// Initialise the database

	db := utils.GetDB()

	query := "SELECT research.research_id, research.research_status, research.research_prototype_url, research.research_intro, research.research_intro_title, research.research_intro_description, questions.question_id, questions.question_title, questions.question_type, questions.question_index, options.* FROM research LEFT JOIN questions ON research.research_id = questions.question_research_id LEFT JOIN options ON research.research_id = options.option_research_id WHERE research.research_id = $1"

	var finalSurvey FinalSurvey

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
		var result JoinedSurvey

		scanErr := rows.Scan(
			&result.ResearchID,
			&result.ResearchStatus,
			&result.ResearchPrototypeUrl,
			&result.ResearchIntro,
			&result.ResearchIntroTitle,
			&result.ResearchIntroDescription,
			&result.QuestionID,
			&result.QuestionTitle,
			&result.QuestionType,
			&result.QuestionIndex,
			&result.OptionID,
			&result.OptionContent,
			&result.OptionQuestionId,
			&result.OptionIndex,
			&result.OptionResearchId,
		)

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

		question := SurveyQuestion{
			QuestionID:    result.QuestionID.String,
			QuestionTitle: result.QuestionTitle.String,
			QuestionType:  result.QuestionType.String,
			QuestionIndex: result.QuestionIndex.Int32,
		}

		option := FetchedOption{
			OptionID:         result.OptionID.String,
			OptionContent:    result.OptionContent.String,
			OptionQuestionId: result.OptionQuestionId.String,
			OptionIndex:      int(result.OptionIndex.Int32),
			OptionResearchId: result.OptionResearchId.String,
		}

		if finalSurvey.ResearchID != result.ResearchID {
			finalSurvey.ResearchID = result.ResearchID
			finalSurvey.ResearchStatus = result.ResearchStatus
			finalSurvey.ResearchPrototypeUrl = result.ResearchPrototypeUrl
			finalSurvey.ResearchIntro = result.ResearchIntro
			finalSurvey.ResearchIntroTitle = result.ResearchIntroTitle
			finalSurvey.ResearchIntroDescription = result.ResearchIntroDescription
		}

		// Check if the question already exists in finalResearch
		questionExists := false
		for _, q := range finalSurvey.ResearchQuestions {
			if q.QuestionID == question.QuestionID {
				questionExists = true
				break
			}
		}

		// If the question doesn't exist, append it to finalResearch
		if !questionExists {
			finalSurvey.ResearchQuestions = append(finalSurvey.ResearchQuestions, question)
		}

		for i, q := range finalSurvey.ResearchQuestions {
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
					finalSurvey.ResearchQuestions[i].QuestionOptions = append(finalSurvey.ResearchQuestions[i].QuestionOptions, option)
				}
			}
		}

	}

	json, _ := json.Marshal(finalSurvey)

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
