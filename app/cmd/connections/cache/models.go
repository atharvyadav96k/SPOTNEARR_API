package cache

import "github.com/redis/go-redis/v9"

type Cache struct {
	passwordSession *passwordSession
}

func newCache(client *redis.Client) *Cache {
	return &Cache{
		passwordSession: newPasswordSession(client),
	}
}
