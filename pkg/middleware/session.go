package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// SessionValidation verifies the ?session= query parameter as a signed JWT.
func SessionValidation(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := strings.TrimSpace(r.URL.Query().Get("session"))
			if session == "" {
				http.Error(w, "Invalid URL", http.StatusBadGateway)
				return
			}
			token, err := jwt.Parse(session, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})
			if err != nil || !token.Valid {
				http.Error(w, "Invalid URL", http.StatusBadGateway)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
