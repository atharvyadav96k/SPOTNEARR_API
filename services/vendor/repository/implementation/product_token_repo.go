package implementation

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type ProductTokenRepository struct{ db *gorm.DB }

func NewProductTokenRepository(db *gorm.DB) *ProductTokenRepository {
	return &ProductTokenRepository{db: db}
}

func (r *ProductTokenRepository) GetTokenCategoryFreqs(ctx context.Context, tokens []string) (map[uint]int, error) {
	if len(tokens) == 0 {
		return map[uint]int{}, nil
	}
	tokenArray := "{" + strings.Join(tokens, ",") + "}"
	type row struct {
		CategoryID uint
		Total      int
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT category_id, SUM(count) AS total
		FROM product_tokens
		WHERE token = ANY(?::text[])
		GROUP BY category_id
	`, tokenArray).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	result := make(map[uint]int, len(rows))
	for _, row := range rows {
		result[row.CategoryID] = row.Total
	}
	return result, nil
}
