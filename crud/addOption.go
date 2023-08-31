package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type Option struct {
	Title string `json:"title"`
	Index int    `json:"index"`
}

func AddOption(w http.ResponseWriter, r *http.Request) {
	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
		return
	}

	var optionArray []Option

	id := r.URL.Query().Get("id")

	err := json.NewDecoder(r.Body).Decode(&optionArray)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonData, err := json.Marshal(optionArray)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db := utils.GetDB()

	query := "UPDATE questions SET options = $1 WHERE id = $2"

	_, updateErr := db.Exec(query, id, jsonData)

	if updateErr != nil {
		fmt.Println(updateErr.Error())
		http.Error(w, "Unable to update columns in database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}
