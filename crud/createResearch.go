package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type Research struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	Limit        int    `json:"limit"`
	PrototypeUrl string `json:"prototype_url"`
}

func CreateResearch(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
		return
	}

	var nr Research

	err := json.NewDecoder(r.Body).Decode(&nr)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "There's been an error decoding the request body. Please try again.", http.StatusInternalServerError)
		return
	}

	db := utils.GetDB()

	query := "INSERT INTO research (title, description, status, \"limit\", prototype_url, user_id) VALUES ($1, $2, $3, $4, $5, $6)"

	// Decode the token to get the user id

	uid, err := utils.DecodeTokenId(tk)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := db.Exec(query, nr.Title, nr.Description, nr.Status, nr.Limit, nr.PrototypeUrl, uid)

	fmt.Println(data)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal("Successfully created the research")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
