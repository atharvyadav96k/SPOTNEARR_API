package middleware

import (
	"net/http"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/utils"
)

func SessionValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := r.URL.Query().Get("session")
		session = strings.TrimSpace(session)
		if session == "" {
			http.Error(w, "Invalid URL", http.StatusBadGateway)
			return
		}
		ok := utils.IsValidSession(session)
		if !ok {
			http.Error(w, "Invalid URL", http.StatusBadGateway)
			return
		}
		next.ServeHTTP(w, r)
	})
}
