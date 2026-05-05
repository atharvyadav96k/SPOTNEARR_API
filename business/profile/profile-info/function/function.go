package business

import (
	"encoding/json"
	"net/http"
)

func BusinessProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Business Profile successful",
	}

	json.NewEncoder(w).Encode(response)
}
