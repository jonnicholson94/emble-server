package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type EditedComment struct {
	Content string `json:"content"`
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")

	var body EditedComment

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	query := "UPDATE comments SET content = $1 WHERE id = $2"

	_, dbErr := db.Exec(query, body.Content, id)

	if dbErr != nil {
		fmt.Println(dbErr.Error())
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal("Successfully updated the comment")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
