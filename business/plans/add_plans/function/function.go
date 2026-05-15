package plans

import (
	"encoding/json"
	"net/http"
)

func Function(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Business plan added successfully",
	}

	json.NewEncoder(w).Encode(response)
}
