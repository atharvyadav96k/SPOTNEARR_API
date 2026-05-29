package cache

func (c *Cache) GetPasswordSessions() *passwordSession {
	return c.passwordSession
}
