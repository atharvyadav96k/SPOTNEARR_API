package spotlight_operation

import (
	"encoding/json"
	"net/http"
)

func Function(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Spotlight retrieved successfully",
	}
	json.NewEncoder(w).Encode(response)
}
