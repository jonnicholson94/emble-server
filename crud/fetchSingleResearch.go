package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchSingleResearch(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")

	db := utils.GetDB()

	query := "SELECT * FROM research WHERE id = $1"

	var research FetchedResearch

	row := db.QueryRow(query, id)

	scanErr := row.Scan(
		&research.ID,
		&research.Title,
		&research.Description,
		&research.Status,
		&research.Limit,
		&research.PrototypeUrl,
		&research.UserId,
	)

	if scanErr != nil {
		fmt.Println(scanErr)
		http.Error(w, scanErr.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(research)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "There's been a problem processing the json", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

}
