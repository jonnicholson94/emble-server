package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type OrderQuestion struct {
	QuestionID    string `json:"question_id"`
	QuestionIndex int    `json:"question_index"`
}

func UpdateQuestionOrder(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		fmt.Println(tokenErr)
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
		return
	}

	var questions []OrderQuestion

	err := json.NewDecoder(r.Body).Decode(&questions)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "There's been an error decoding the request body. Please try again.", http.StatusInternalServerError)
		return
	}

	db := utils.GetDB()

	for _, data := range questions {
		updateQuery := "UPDATE questions SET index = $1 WHERE id = $2"

		_, err = db.Exec(updateQuery, data.QuestionIndex, data.QuestionID)

		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)

}
