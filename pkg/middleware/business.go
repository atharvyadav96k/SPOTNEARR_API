package middleware

import (
	"net/http"
	"strings"

	"github.com/atharvyadav96k/spotnearr/pkg/jwtutil"
)

// BusinessOnly allows only authenticated users who hold a business role.
// It re-parses the token if claims aren't already in context (e.g. standalone use).
func BusinessOnly(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r)
			if !ok {
				parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					parsed, err := jwtutil.ValidateToken(parts[1], jwtSecret, jwtutil.TypeAccessToken)
					if err == nil {
						claims = parsed
						ok = true
					}
				}
			}
			if !ok || claims.BusinessId == nil || *claims.BusinessId == 0 || claims.UserRole == jwtutil.RoleUser {
				http.Error(w, "Unauthorized access", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
