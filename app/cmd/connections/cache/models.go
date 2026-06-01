package cache

import "github.com/redis/go-redis/v9"

type Cache struct {
	passwordSession       *passwordSession
	refresh_token_session *refresh_token_session
}

func newCache(client *redis.Client) *Cache {
	return &Cache{
		passwordSession:       newPasswordSession(client),
		refresh_token_session: newRefreshTokenSession(client),
	}
}
