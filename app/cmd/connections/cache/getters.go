package cache

import pkgrl "github.com/atharvyadav96k/spotnearr/pkg/ratelimit"

func (c *Cache) GetPasswordSessions() *passwordSession {
	return c.passwordSession
}

func (c *Cache) GetRefreshTokenSession() *refresh_token_session {
	return c.refresh_token_session
}

func (c *Cache) GetRateLimit() *pkgrl.RateLimiter {
	return c.rateLimiter
}

func (c *Cache) GetDealGeoCache() *dealGeoCache {
	return c.dealGeoCache
}
