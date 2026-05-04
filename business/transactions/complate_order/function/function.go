package transactions

import (
	"encoding/json"
	"net/http"
)

func CompleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]string{
		"message": "Order completion successful",
	}

	json.NewEncoder(w).Encode(response)
}
