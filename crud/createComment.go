package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewComment struct {
	Content    string `json:"content"`
	ResearchId string `json:"research_id"`
	Timestamp  int    `json:"timestamp"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, "User's token is invalid", http.StatusUnauthorized)
		return
	}

	uid, err := utils.DecodeTokenId(tk)

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var nc NewComment

	decodeErr := json.NewDecoder(r.Body).Decode(&nc)

	if decodeErr != nil {
		fmt.Println(decodeErr)
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}

	db := utils.GetDB()

	insertQuery := "INSERT INTO comments (content, user_id, research_id, timestamp) VALUES ($1, $2, $3, $4)"

	fmt.Println(nc)

	data, err := db.Exec(insertQuery, nc.Content, uid, nc.ResearchId, nc.Timestamp)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to save data, please try again", http.StatusInternalServerError)
		return
	}

	fmt.Println(data)

	res, err := json.Marshal("Successfully created the comment")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
