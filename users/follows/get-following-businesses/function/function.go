package follow

import (
	"encoding/json"
	"net/http"
)

func GetFollowedBusiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Followed Businesses retrieval successful",
	}

	json.NewEncoder(w).Encode(response)
}
