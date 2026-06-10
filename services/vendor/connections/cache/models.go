package cache

import "github.com/redis/go-redis/v9"

type Cache struct {
	client     *redis.Client
	rate_limit *rate_limit
}

func newCache(client *redis.Client) *Cache {
	return &Cache{
		client:     client,
		rate_limit: newRateLimit(client),
	}
}

func (c *Cache) RedisClient() *redis.Client { return c.client }
