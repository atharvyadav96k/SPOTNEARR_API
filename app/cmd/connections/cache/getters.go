package cache

func (c *Cache) GetPasswordSessions() *passwordSession {
	return c.passwordSession
}

func (c *Cache) GetRefreshTokenSession() *refresh_token_session {
	return c.refresh_token_session
}

func (c *Cache) GetRateLimit() *rate_limit {
	return c.rate_limit
}
