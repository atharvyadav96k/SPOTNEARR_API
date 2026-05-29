package cache

import "github.com/redis/go-redis/v9"

type base_cache struct {
	client *redis.Client
}

func newBaseCache(client *redis.Client) *base_cache {
	return &base_cache{
		client: client,
	}
}

func (b *base_cache) getClient() *redis.Client {
	return b.client
}
