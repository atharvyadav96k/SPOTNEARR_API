package spotlight_interaction

import (
	"encoding/json"
	"net/http"
)

func SpotlightLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Liked successfully",
	}

	json.NewEncoder(w).Encode(response)
}
