package crud

import (
	"emble-server/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

// 1. Check token
// 2. Prepare database
// 3. Write query
// 4. Look into cascading questions and research so they're deleted together
// 5. Return success

func DeleteResearch(w http.ResponseWriter, r *http.Request) {

	tk := r.Header.Get("Authorization")

	tokenErr := utils.ValidateToken(tk)

	if tokenErr != nil {
		http.Error(w, "User's token is invalid", http.StatusUnauthorized)
		return
	}

	db := utils.GetDB()

	id := r.URL.Query().Get("id")

	_, queErr := db.Exec("DELETE FROM questions WHERE research_id = $1", id)

	if queErr != nil {
		fmt.Println(queErr.Error())
		http.Error(w, "Failed to delete the questions associated with the research", http.StatusBadRequest)
		return
	}

	_, resErr := db.Exec("DELETE FROM research WHERE id = $1", id)

	if resErr != nil {
		fmt.Println(resErr.Error())
		http.Error(w, "Failed to delete the research", http.StatusBadRequest)
		return
	}

	res, jsonErr := json.Marshal("Successfully deleted the research")

	if jsonErr != nil {
		http.Error(w, "Failed to marshal the json to return to the frontend", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(res)

}
