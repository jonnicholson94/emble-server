package crud

import (
	"emble-server/utils"
	"encoding/json"
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
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
		return
	}

	var nq NewQuestion

	err := json.NewDecoder(r.Body).Decode(&nq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	insertQuery := "INSERT INTO questions (title, type, research_id, index) VALUES ($1, $2, $3, $4)"

	data, err := db.Exec(insertQuery, nq.Title, nq.Type, nq.ResearchId, nq.Index)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(data)

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
