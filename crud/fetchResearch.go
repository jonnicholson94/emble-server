package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type FetchedResearch struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Limit        string `json:"limit"`
	PrototypeUrl string `json:"prototype_url"`
	UserId       int    `json:"user_id"`
}

func FetchResearch(w http.ResponseWriter, r *http.Request) {

	token := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(token)

	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
		return
	}

	uid, err := utils.DecodeTokenId(token)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	query := "SELECT * FROM research WHERE research_user_id = $1"

	rows, err := db.Query(query, uid)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var results []FetchedResearch

	for rows.Next() {
		var result FetchedResearch

		err := rows.Scan(
			&result.ID,
			&result.Title,
			&result.Description,
			&result.Status,
			&result.Limit,
			&result.PrototypeUrl,
			&result.UserId,
		)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		results = append(results, result)
	}

	data, err := json.Marshal(results)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
