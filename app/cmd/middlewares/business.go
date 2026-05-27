package middleware

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/auth"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

func BusinessOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*auth.UserClaims)
		if !ok {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}
		if claims.BusinessId == nil || *claims.BusinessId == 0 || claims.UserRole == models.RoleUser {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
