package crud

import (
	"emble-server/utils"
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
		http.Error(w, tokenErr.Error(), http.StatusUnauthorized)
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

	w.WriteHeader(http.StatusOK)

}
