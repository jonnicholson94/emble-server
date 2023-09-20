package questions

import (
	"emble-server/auth"
	"emble-server/utils"
	"encoding/json"
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
		updateQuery := "UPDATE questions SET index = $1 WHERE id = $2"

		_, err = db.Exec(updateQuery, data.QuestionIndex, data.QuestionID)

		if err != nil {
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

	w.WriteHeader(http.StatusOK)

}
