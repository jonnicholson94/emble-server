package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func EditItem(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
		return
	}

	body := make(map[string]interface{})

	id := r.URL.Query().Get("id")

	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	updateQuery := fmt.Sprintf("UPDATE research SET %s WHERE id = $%d", strings.Join(updateColumns, ", "), i)

	values = append(values, id)

	_, err = db.Exec(updateQuery, values...)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
}
