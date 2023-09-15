package crud

import (
	"emble-server/utils"
	"encoding/json"
	"net/http"
)

type EditedComment struct {
	CommentContent string `json:"comment_content"`
}

func EditComment(w http.ResponseWriter, r *http.Request) {
	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

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

	id := r.URL.Query().Get("id")

	var body EditedComment

	err := json.NewDecoder(r.Body).Decode(&body)

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

	query := "UPDATE comments SET comment_content = $1 WHERE comment_id = $2"

	_, dbErr := db.Exec(query, body.CommentContent, id)

	if dbErr != nil {
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

	res, err := json.Marshal("Successfully updated the comment")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
