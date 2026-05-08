package user_spotlights

import (
	"encoding/json"
	"net/http"
)

func Function(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Spotlights retrieval successfully",
	}

	json.NewEncoder(w).Encode(response)
}
