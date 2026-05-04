package spotlight_operation

import (
	"encoding/json"
	"net/http"
)

func GetSpotLights(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Login successful",
	}
	json.NewEncoder(w).Encode(response)
}
