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

type CustomError struct {
	Message string `json:"message"`
	Status  int32  `json:"status"`
}

func CreateComment(w http.ResponseWriter, r *http.Request) {

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

	uid, err := utils.DecodeTokenId(tk)

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

	insertQuery := "INSERT INTO comments (content, user_id, research_id, timestamp) VALUES ($1, $2, $3, $4)"

	fmt.Println(nc)

	data, err := db.Exec(insertQuery, nc.Content, uid, nc.ResearchId, nc.Timestamp)

	if err != nil {
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
