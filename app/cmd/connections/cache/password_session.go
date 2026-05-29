package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type passwordSession struct {
	base_cache
}

func newPasswordSession(client *redis.Client) *passwordSession {
	return &passwordSession{
		base_cache: *newBaseCache(client),
	}
}

func (p *passwordSession) key(key string) string {
	return fmt.Sprint("passwordSession:", key)
}

func (p *passwordSession) NewPasswordSession(ctx context.Context, email string, expireTime time.Duration, session string) (string, error) {
	val := p.getClient().Get(ctx, p.key(email)).Val()
	if val != "" {
		return val, nil
	}
	cache := p.getClient().Set(ctx, p.key(session), email, time.Duration(expireTime))
	if cache.Err() != nil {
		return "", cache.Err()
	}
	return cache.Val(), nil
}

func (p *passwordSession) GetPasswordSession(ctx context.Context, session string) (string, error) {
	cache := p.getClient().GetDel(ctx, p.key(session))
	if cache.Err() != nil {
		return "", nil
	}
	return cache.Val(), nil
}
