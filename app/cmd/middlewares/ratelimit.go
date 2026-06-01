package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/atharvyadav96k/SPOTNEARR_API/auth"
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
)

func RateLimit(c *cache.Cache, limit int, window time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(ClaimsKey).(*auth.UserClaims)
			if !ok || claims == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			allowed, err := c.GetRateLimit().Allow(
				context.Background(),
				claims.UserId,
				r.URL.Path,
				limit,
				window,
			)
			if err != nil {
				log.Default().Println("Rate limit check failed:", err)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				http.Error(w, "Too many requests, please try again later", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
