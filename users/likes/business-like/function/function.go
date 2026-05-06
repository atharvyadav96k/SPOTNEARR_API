package user_likes

import (
	"encoding/json"
	"net/http"
)

func LikeBusiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Business liked successfully",
	}
	json.NewEncoder(w).Encode(response)
}
