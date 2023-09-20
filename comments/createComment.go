package comments

import (
	"emble-server/auth"
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewComment struct {
	CommentId         string `json:"comment_id"`
	CommentContent    string `json:"comment_content"`
	CommentResearchId string `json:"comment_research_id"`
	CommentTimestamp  int    `json:"comment_timestamp"`
}

type CustomError struct {
	Message string `json:"message"`
	Status  int32  `json:"status"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {

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

	uid, err := auth.DecodeTokenId(tk)

	if err != nil {
		customErr := CustomError{
			Message: "Error decoding token",
			Status:  http.StatusUnauthorized,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	var nc NewComment

	decodeErr := json.NewDecoder(r.Body).Decode(&nc)

	if decodeErr != nil {
		customErr := CustomError{
			Message: "Error decoding body",
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

	insertQuery := "INSERT INTO comments (comment_id, comment_content, comment_user_id, comment_research_id, comment_timestamp) VALUES ($1, $2, $3, $4, $5)"

	fmt.Println(nc)

	data, err := db.Exec(insertQuery, nc.CommentId, nc.CommentContent, uid, nc.CommentResearchId, nc.CommentTimestamp)

	if err != nil {
		fmt.Println(err)
		customErr := CustomError{
			Message: "Database error",
			Status:  http.StatusInternalServerError,
		}

		// Convert the error to JSON
		errJSON, _ := json.Marshal(customErr)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errJSON)
		return
	}

	fmt.Println(data)

	res, _ := json.Marshal("Successfully created the comment")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
