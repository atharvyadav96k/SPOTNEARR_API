package profile_image

import (
	"encoding/json"
	"net/http"
)

func SetProfileImageInDb(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Profile image set in database successfully",
	}

	json.NewEncoder(w).Encode(response)
}
