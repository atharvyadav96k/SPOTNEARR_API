package ratelimit

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

// slidingWindow atomically enforces a sliding-log rate limit using a sorted set.
// Score = Unix ms timestamp; member = timestamp+rand so same-ms requests don't collapse.
// Returns 1 if allowed, 0 if limit exceeded.
var slidingWindow = redis.NewScript(`
local key      = KEYS[1]
local now      = tonumber(ARGV[1])
local windowMs = tonumber(ARGV[2])
local limit    = tonumber(ARGV[3])
local member   = ARGV[4]
local cutoff   = now - windowMs

redis.call('ZREMRANGEBYSCORE', key, '-inf', cutoff)
local count = redis.call('ZCARD', key)
redis.call('PEXPIRE', key, windowMs)
if count < limit then
    redis.call('ZADD', key, now, member)
    return 1
end
return 0
`)

// RateLimiter is a Redis-backed sliding window rate limiter.
type RateLimiter struct {
	client    *redis.Client
	keyPrefix string
}

// New returns a RateLimiter that prefixes all keys with keyPrefix.
// Pass the APP_ENV value (e.g. "prod:") as keyPrefix for namespace isolation.
func New(client *redis.Client, keyPrefix string) *RateLimiter {
	return &RateLimiter{client: client, keyPrefix: keyPrefix}
}

// Allow returns true if the request from userID on route is within the limit
// for the given sliding window duration.
func (r *RateLimiter) Allow(ctx context.Context, userID any, route string, limit int, window time.Duration) (bool, error) {
	key := r.keyPrefix + fmt.Sprintf("rateLimit:%v:%s", userID, route)
	nowMs := time.Now().UnixMilli()
	member := fmt.Sprintf("%d-%d", nowMs, rand.Int63())

	result, err := slidingWindow.Run(ctx, r.client,
		[]string{key}, nowMs, window.Milliseconds(), limit, member,
	).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
