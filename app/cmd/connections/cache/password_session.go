package cache

import (
	"context"
	"fmt"
	"os"
	"strconv"
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

func (p *passwordSession) NewPasswordSession(ctx context.Context, email string, session string) error {
	val := p.getClient().Get(ctx, p.key(email)).Val()
	if val != "" {
		return nil
	}
	passwordSessionTime := 1
	if val := os.Getenv("PASSWORD_SESSION_TIME"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			passwordSessionTime = parsed
		}
	}
	return p.getClient().Set(ctx, p.key(email), session, time.Duration(passwordSessionTime)*time.Hour).Err()
}
