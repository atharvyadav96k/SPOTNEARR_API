package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/atharvyadav96k/spotnearr/pkg/jwtutil"
)

type contextKey string

const ClaimsKey contextKey = "userClaims"

// ClaimsFromContext extracts validated UserClaims from the request context.
func ClaimsFromContext(r *http.Request) (*jwtutil.UserClaims, bool) {
	claims, ok := r.Context().Value(ClaimsKey).(*jwtutil.UserClaims)
	return claims, ok && claims != nil
}

// Auth validates a Bearer JWT access token and stores the claims in the context.
func Auth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Default().Println("Missing auth header")
				http.Error(w, "Unauthorized access", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				http.Error(w, "Unauthorized access", http.StatusUnauthorized)
				return
			}
			claims, err := jwtutil.ValidateToken(parts[1], jwtSecret, jwtutil.TypeAccessToken)
			if err != nil {
				log.Default().Println("Failed to validate token:", err)
				http.Error(w, "Unauthorized access", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
