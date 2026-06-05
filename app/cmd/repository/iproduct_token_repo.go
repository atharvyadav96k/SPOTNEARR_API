package repository

import "context"

type IProductTokenRepository interface {
	GetTokenCategoryFreqs(ctx context.Context, tokens []string) (map[uint]int, error)
}
