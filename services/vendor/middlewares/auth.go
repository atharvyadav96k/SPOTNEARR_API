package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/atharvyadav96k/spotnearr/pkg/jwtutil"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/config"
)

type contextKey string

const ClaimsKey contextKey = "userClaims"

func Auth(next http.Handler) http.Handler {
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
		claims, err := jwtutil.ValidateToken(parts[1], config.C.JWTSecret, jwtutil.TypeAccessToken)
		if err != nil {
			log.Default().Println("Failed to validate token:", err)
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}
		if claims.TokenType != jwtutil.TypeAccessToken {
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
