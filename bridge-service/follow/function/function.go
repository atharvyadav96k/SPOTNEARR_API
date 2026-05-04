package follow_action

import (
	"encoding/json"
	"net/http"
)

func FollowAction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Follow action executed successfully",
	}

	json.NewEncoder(w).Encode(response)
}
