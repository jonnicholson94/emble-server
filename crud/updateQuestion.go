package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func UpdateQuestion(w http.ResponseWriter, r *http.Request) {
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

	body := make(map[string]interface{})

	id := r.URL.Query().Get("id")

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

	// Construct the dynamic query
	var updateColumns []string
	var values []interface{}
	i := 1
	for key, value := range body {
		updateColumns = append(updateColumns, fmt.Sprintf("%s = $%d", key, i))
		values = append(values, value)
		i++
	}

	// Construct and execute the query
	updateQuery := fmt.Sprintf("UPDATE questions SET %s WHERE id = $%d", strings.Join(updateColumns, ", "), i)

	values = append(values, id)

	_, err = db.Exec(updateQuery, values...)
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

	res, err := json.Marshal("Successfully saved your changes")

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
