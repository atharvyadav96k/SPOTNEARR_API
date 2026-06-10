package search

import "context"

// SearchRow is the scan target for search queries — includes the computed
// token_match_cnt from the LATERAL subquery (not a real column in search_entries).
type SearchRow struct {
	SearchEntry
	TokenMatchCnt int `gorm:"column:token_match_cnt"`
}

type ISearchRepository interface {
	Search(ctx context.Context, tokens []string, lat, long *float64, rangeKm float64) ([]SearchRow, error)
	Upsert(ctx context.Context, entry SearchEntry) error
	SoftDelete(ctx context.Context, id uint) error
}
