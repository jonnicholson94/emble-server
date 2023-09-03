package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewQuestion struct {
	Title      string `json:"title"`
	Type       string `json:"type"`
	ResearchId string `json:"research_id"`
	Index      int    `json:"index"`
}

func CreateQuestion(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, "User's token is invalid", http.StatusUnauthorized)
		return
	}

	var nq NewQuestion

	err := json.NewDecoder(r.Body).Decode(&nq)

	if err != nil {
		http.Error(w, "Failed to process the request, please try again", http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	insertQuery := "INSERT INTO questions (title, type, research_id, index) VALUES ($1, $2, $3, $4)"

	data, err := db.Exec(insertQuery, nq.Title, nq.Type, nq.ResearchId, nq.Index)

	fmt.Println(data)

	if err != nil {
		http.Error(w, "Failed to save data, please try again", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal("Successfully created the question")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
