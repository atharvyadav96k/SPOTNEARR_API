package plans

import (
	"encoding/json"
	"net/http"
)

func SetActivePlan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "",
	}

	json.NewEncoder(w).Encode(response)
}
