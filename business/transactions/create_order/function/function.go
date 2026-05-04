package transactions

import (
	"encoding/json"
	"net/http"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Order creation successful",
	}

	json.NewEncoder(w).Encode(response)
}
