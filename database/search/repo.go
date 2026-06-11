package search

import "context"

// SearchRow is the scan target for search queries — includes the computed
// token_match_cnt from the LATERAL subquery (not a real column in search_entries).
type SearchRow struct {
	SearchEntry
	TokenMatchCnt int `gorm:"column:token_match_cnt"`
}

// FreqDelta is a signed change to apply to a token-category frequency counter.
// Positive = token added to a product in that category. Negative = removed.
type FreqDelta struct {
	Token      string
	CategoryID uint
	Delta      int64
}

type ISearchRepository interface {
	Search(ctx context.Context, tokens []string, lat, long *float64, rangeKm float64) ([]SearchRow, error)
	GetByID(ctx context.Context, id uint) (*SearchEntry, error)
	Upsert(ctx context.Context, entry SearchEntry) error
	SoftDelete(ctx context.Context, id uint) error
	GetTokenCategoryFreqs(ctx context.Context, tokens []string) ([]TokenCategoryFreq, error)
	FlushFreqDeltas(ctx context.Context, deltas []FreqDelta) error
}
