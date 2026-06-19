package cache

import (
	"os"

	pkgrl "github.com/atharvyadav96k/spotnearr/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	passwordSession       *passwordSession
	refresh_token_session *refresh_token_session
	rateLimiter           *pkgrl.RateLimiter
	dealGeoCache          *dealGeoCache
}

func newCache(client *redis.Client) *Cache {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "prod"
	}
	return &Cache{
		passwordSession:       newPasswordSession(client),
		refresh_token_session: newRefreshTokenSession(client),
		rateLimiter:           pkgrl.New(client, env+":"),
		dealGeoCache:          newDealGeoCache(client),
	}
}
