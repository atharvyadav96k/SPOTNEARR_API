package middleware

import (
	"context"
	"log"
	"net/http"
	"time"
)

// RateLimiter is satisfied by any Redis-backed rate limit implementation.
type RateLimiter interface {
	Allow(ctx context.Context, userID any, route string, limit int, window time.Duration) (bool, error)
}

// RateLimit gates requests per authenticated user per path.
// Pass cache.GetRateLimit() as the limiter from any service's cache package.
func RateLimit(limiter RateLimiter, limit int, window time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			allowed, err := limiter.Allow(context.Background(), claims.UserId, r.URL.Path, limit, window)
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
