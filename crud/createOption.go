package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type Option struct {
	Content    string `json:"content"`
	QuestionID string `json:"question_id"`
	Index      int    `json:"index"`
	ResearchID string `json:"research_id"`
}

func CreateOption(w http.ResponseWriter, r *http.Request) {
	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
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

	query := "INSERT INTO options (content, question_id, index, research_id) VALUES ($1, $2, $3, $4)"

	_, insertErr := db.Exec(query, option.Content, option.QuestionID, option.Index, option.ResearchID)

	if insertErr != nil {
		fmt.Println(insertErr.Error())
		http.Error(w, insertErr.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal("Successfully added the option")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
