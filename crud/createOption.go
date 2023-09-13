package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type Option struct {
	OptionContent    string `json:"option_content"`
	OptionQuestionID string `json:"option_question_id"`
	OptionIndex      int    `json:"option_index"`
	OptionResearchID string `json:"option_research_id"`
}

func CreateOption(w http.ResponseWriter, r *http.Request) {
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

	var option Option

	err := json.NewDecoder(r.Body).Decode(&option)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	query := "INSERT INTO options (option_content, option_question_id, option_index, option_research_id) VALUES ($1, $2, $3, $4)"

	_, insertErr := db.Exec(query, option.OptionContent, option.OptionQuestionID, option.OptionIndex, option.OptionResearchID)

	if insertErr != nil {
		customErr := CustomError{
			Message: "Failed to insert the option",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	res, err := json.Marshal("Successfully added the option")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
