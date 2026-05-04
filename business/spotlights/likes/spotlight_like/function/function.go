package spotlight_interaction

import (
	"encoding/json"
	"net/http"
)

func SpotlightDisLike(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Disliked successfully",
	}

	json.NewEncoder(w).Encode(response)
}
