package plans

import (
	"encoding/json"
	"net/http"
)

func Function(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Active plan",
	}

	json.NewEncoder(w).Encode(response)
}
