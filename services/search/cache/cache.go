package cache

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	searchdb "github.com/Developer-Aadesh/spotnearr-database/search"
	"github.com/redis/go-redis/v9"
)

const freqCacheTTL = 30 * time.Second

// SearchCache wraps Redis operations used by the search service.
type SearchCache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *SearchCache {
	return &SearchCache{rdb: rdb}
}

// GetTokenFreqs returns cached token-category frequencies for the given tokens.
// Returns nil, false on a cache miss or any Redis error.
func (c *SearchCache) GetTokenFreqs(ctx context.Context, tokens []string) ([]searchdb.TokenCategoryFreq, bool) {
	b, err := c.rdb.Get(ctx, freqKey(tokens)).Bytes()
	if err != nil {
		return nil, false
	}
	var freqs []searchdb.TokenCategoryFreq
	if json.Unmarshal(b, &freqs) != nil {
		return nil, false
	}
	return freqs, true
}

// SetTokenFreqs stores token-category frequencies in Redis. Errors are silently
// ignored — a Redis outage must never break search.
func (c *SearchCache) SetTokenFreqs(ctx context.Context, tokens []string, freqs []searchdb.TokenCategoryFreq) {
	b, err := json.Marshal(freqs)
	if err != nil {
		return
	}
	c.rdb.Set(ctx, freqKey(tokens), b, freqCacheTTL)
}

// freqKey builds a stable cache key from a set of tokens.
// Tokens are sorted so "apple juice" and "juice apple" share the same key.
func freqKey(tokens []string) string {
	sorted := make([]string, len(tokens))
	copy(sorted, tokens)
	sort.Strings(sorted)
	return "search:freqs:" + strings.Join(sorted, ",")
}
