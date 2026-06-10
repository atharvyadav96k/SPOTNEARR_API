package cache

import (
	"os"

	pkgrl "github.com/atharvyadav96k/spotnearr/pkg/ratelimit"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client      *redis.Client
	rateLimiter *pkgrl.RateLimiter
}

func newCache(client *redis.Client) *Cache {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "prod"
	}
	return &Cache{
		client:      client,
		rateLimiter: pkgrl.New(client, env+":"),
	}
}

func (c *Cache) RedisClient() *redis.Client { return c.client }
