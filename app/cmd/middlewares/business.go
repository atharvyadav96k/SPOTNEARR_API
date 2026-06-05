package middleware

import (
	"net/http"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/auth"
	"github.com/atharvyadav96k/SPOTNEARR_API/config"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

func BusinessOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*auth.UserClaims)

		// Auth middleware sets claims in context; fall back to parsing the
		// Authorization header directly in case context propagation failed.
		if !ok || claims == nil {
			parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				parsed, err := auth.ValidateToken(parts[1], config.C.JWTSecret, auth.TypeAccessToken)
				if err == nil {
					claims = parsed
					ok = true
				}
			}
		}

		if !ok || claims == nil {
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
