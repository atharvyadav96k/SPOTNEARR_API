package cache

func (c *Cache) GetRateLimit() *rate_limit {
	return c.rate_limit
}
