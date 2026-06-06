package cache

import (
	"os"

	"github.com/redis/go-redis/v9"
)

type base_cache struct {
	client *redis.Client
	env    string
}

func newBaseCache(client *redis.Client) *base_cache {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "prod"
	}
	return &base_cache{
		client: client,
		env:    env,
	}
}

func (b *base_cache) getClient() *redis.Client {
	return b.client
}

func (b *base_cache) pk(key string) string {
	return b.env + ":" + key
}
