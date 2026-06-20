package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	tcKeyPrefix  = "search:tc:"   // sorted set: token → {catID: freq}
	docKeyPrefix = "search:doc:"  // string: invProductID → {tokens, catIDs}
	docTTL       = 7 * 24 * time.Hour
)

type docMeta struct {
	Tokens  []string `json:"t"`
	CatIDs  []int64  `json:"c"`
}

type TokenCategoryCache struct {
	client *redis.Client
}

func New(client *redis.Client) *TokenCategoryCache {
	return &TokenCategoryCache{client: client}
}

// IncrFreqs increments token×category frequency counts for a newly indexed product.
// Uses a single pipeline — one round-trip for all token×category pairs.
func (c *TokenCategoryCache) IncrFreqs(ctx context.Context, tokens []string, catIDs []int64) error {
	if len(tokens) == 0 || len(catIDs) == 0 {
		return nil
	}
	pipe := c.client.Pipeline()
	for _, tok := range tokens {
		key := tcKeyPrefix + tok
		for _, catID := range catIDs {
			pipe.ZIncrBy(ctx, key, 1, strconv.FormatInt(catID, 10))
		}
	}
	_, err := pipe.Exec(ctx)
	return err
}

// DecrFreqs decrements token×category frequency counts when a product is deleted.
func (c *TokenCategoryCache) DecrFreqs(ctx context.Context, tokens []string, catIDs []int64) error {
	if len(tokens) == 0 || len(catIDs) == 0 {
		return nil
	}
	pipe := c.client.Pipeline()
	for _, tok := range tokens {
		key := tcKeyPrefix + tok
		for _, catID := range catIDs {
			pipe.ZIncrBy(ctx, key, -1, strconv.FormatInt(catID, 10))
		}
	}
	_, err := pipe.Exec(ctx)
	return err
}

// GetTopCategories looks up the highest-frequency categories across all query tokens
// and returns up to topN deduplicated category IDs. Returns nil on any Redis error
// so the caller can degrade gracefully (skip the _eval boost).
func (c *TokenCategoryCache) GetTopCategories(ctx context.Context, tokens []string, topN int) []int64 {
	if len(tokens) == 0 {
		return nil
	}

	// vote[catID] = number of tokens that rank this category in their top-N
	vote := make(map[int64]int, topN*len(tokens))

	pipe := c.client.Pipeline()
	cmds := make([]*redis.StringSliceCmd, len(tokens))
	for i, tok := range tokens {
		cmds[i] = pipe.ZRevRange(ctx, tcKeyPrefix+tok, 0, int64(topN-1))
	}
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		log.Printf("search: GetTopCategories redis: %v", err)
		return nil
	}

	for _, cmd := range cmds {
		for _, member := range cmd.Val() {
			if id, err := strconv.ParseInt(member, 10, 64); err == nil {
				vote[id]++
			}
		}
	}

	if len(vote) == 0 {
		return nil
	}

	// Collect unique IDs ordered by vote count (simple selection — topN is small)
	result := make([]int64, 0, topN)
	for len(result) < topN && len(vote) > 0 {
		var best int64
		var bestVote int
		for id, v := range vote {
			if v > bestVote {
				bestVote = v
				best = id
			}
		}
		result = append(result, best)
		delete(vote, best)
	}
	return result
}

// SetDocMeta stores tokens and category IDs for a product so they can be
// retrieved when the product is deleted (to decrement frequencies).
func (c *TokenCategoryCache) SetDocMeta(ctx context.Context, invProductID uint, tokens []string, catIDs []int64) error {
	key := fmt.Sprintf("%s%d", docKeyPrefix, invProductID)
	b, err := json.Marshal(docMeta{Tokens: tokens, CatIDs: catIDs})
	if err != nil {
		return err
	}
	return c.client.Set(ctx, key, b, docTTL).Err()
}

// GetDocMeta retrieves stored tokens and category IDs for a product.
// Returns nil slices if the key is missing (TTL expired or never set).
func (c *TokenCategoryCache) GetDocMeta(ctx context.Context, invProductID uint) ([]string, []int64) {
	key := fmt.Sprintf("%s%d", docKeyPrefix, invProductID)
	b, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, nil
	}
	var m docMeta
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, nil
	}
	return m.Tokens, m.CatIDs
}

// DelDocMeta removes the stored metadata for a product after deletion.
func (c *TokenCategoryCache) DelDocMeta(ctx context.Context, invProductID uint) {
	c.client.Del(ctx, fmt.Sprintf("%s%d", docKeyPrefix, invProductID))
}
