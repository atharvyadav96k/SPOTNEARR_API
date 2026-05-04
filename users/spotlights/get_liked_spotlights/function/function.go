package spotlights

import (
	"encoding/json"
	"net/http"
)

func SpotlightLikedBusiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Liked Spotlights retrieval successful",
	}

	json.NewEncoder(w).Encode(response)
}
