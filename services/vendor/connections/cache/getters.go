package cache

import pkgrl "github.com/atharvyadav96k/spotnearr/pkg/ratelimit"

func (c *Cache) GetRateLimit() *pkgrl.RateLimiter {
	return c.rateLimiter
}
