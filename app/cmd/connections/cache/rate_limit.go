// cache/rate_limit.go
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type rate_limit struct {
	base_cache
}

func newRateLimit(client *redis.Client) *rate_limit {
	return &rate_limit{
		base_cache: *newBaseCache(client),
	}
}

func (r *rate_limit) key(userID any, route string) string {
	return fmt.Sprintf("rateLimit:%v:%s", userID, route)
}

func (r *rate_limit) Allow(ctx context.Context, userID any, route string, limit int, window time.Duration) (bool, error) {
	key := r.pk(r.key(userID, route))

	count, err := r.getClient().Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if count == 1 {
		r.getClient().Expire(ctx, key, window)
	}

	return count <= int64(limit), nil
}
