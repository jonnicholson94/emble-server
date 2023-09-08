package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func DeleteOption(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, "User's token is invalid", http.StatusUnauthorized)
		return
	}

	db := utils.GetDB()

	id := r.URL.Query().Get("id")

	_, err := db.Exec("DELETE FROM options WHERE option_id = $1", id)

	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Failed to delete the questions associated with the research", http.StatusBadRequest)
		return
	}

	res, jsonErr := json.Marshal("Successfully deleted the option")

	if jsonErr != nil {
		http.Error(w, "Failed to marshal the json to return to the frontend", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
