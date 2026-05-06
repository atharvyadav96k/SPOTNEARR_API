package user_likes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func DiskLikeBusiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": fmt.Sprintf("Business disliked successfully"),
	}

	json.NewEncoder(w).Encode(response)
}
