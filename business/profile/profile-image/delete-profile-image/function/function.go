package profile_image

import (
	"encoding/json"
	"net/http"
)

func DeleteProfileImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Profile image deleted successfully",
	}

	json.NewEncoder(w).Encode(response)
}
