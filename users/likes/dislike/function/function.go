package user_likes

import (
	"encoding/json"
	"net/http"
)

func DisLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Login successful",
	}

	json.NewEncoder(w).Encode(response)
}
