package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type refresh_token_session struct {
	base_cache
}

func newRefreshTokenSession(client *redis.Client) *refresh_token_session {
	return &refresh_token_session{
		base_cache: *newBaseCache(client),
	}
}

func (r *refresh_token_session) key(key any) string {
	return fmt.Sprint("refreshToken:", key)
}

func (r *refresh_token_session) NewRefreshToken(key any, token string, expireTime time.Duration) error {
	cache := r.getClient().Set(context.Background(), r.key(key), token, expireTime)
	if cache.Err() != nil {
		return cache.Err()
	}
	return nil
}

func (r *refresh_token_session) InvalidateRefreshToken(key any) error {
	cache := r.getClient().Del(context.Background(), r.key(key))
	if cache.Err() != nil {
		return cache.Err()
	}
	return nil
}

func (r *refresh_token_session) GetRefreshTokenSession(key any) (string, error) {
	cache := r.getClient().Get(context.Background(), r.key(key))
	if cache.Err() != nil {
		return "", cache.Err()
	}
	return cache.Val(), nil
}
