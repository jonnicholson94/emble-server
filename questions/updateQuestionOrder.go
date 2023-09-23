package questions

import (
	"emble-server/auth"
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

	tokenErr := auth.ValidateToken(tk)

	if tokenErr != nil {
		customErr := CustomError{
			Message: "Invalid token",
			Status:  http.StatusUnauthorized,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errJSON)
		return
	}

	var questions []OrderQuestion

	err := json.NewDecoder(r.Body).Decode(&questions)

	if err != nil {
		fmt.Println(err.Error())
		customErr := CustomError{
			Message: "Failed to process request",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	db := utils.GetDB()

	for _, data := range questions {
		updateQuery := "UPDATE questions SET question_index = $1 WHERE question_id = $2"

		_, err = db.Exec(updateQuery, data.QuestionIndex, data.QuestionID)

		if err != nil {
			fmt.Println(err.Error())
			customErr := CustomError{
				Message: "Failed to process request",
				Status:  http.StatusInternalServerError,
			}

			// Convert the error to JSON
			errJSON, _ := json.Marshal(customErr)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(errJSON)
			return
		}
	}

	res, err := json.Marshal("Successfully saved your changes")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
