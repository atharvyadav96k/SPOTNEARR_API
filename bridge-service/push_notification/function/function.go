package login_function

import (
	"encoding/json"
	"net/http"
)

func PushNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Notification sent successfully",
	}

	json.NewEncoder(w).Encode(response)
}
