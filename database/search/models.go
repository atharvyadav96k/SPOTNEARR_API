package search

import "time"

// SearchEntry is the minimal search index row. One row per inventory_product.
// Stores only what is needed for fast discovery: product reference, tokens, location, categories.
// All business metadata (name, price, description, store info) is deliberately excluded.
type SearchEntry struct {
	ID           uint       `gorm:"primaryKey"`              // = inventory_product.id
	ProductID    uint       `gorm:"index;not null"`
	SearchTokens []string   `gorm:"type:jsonb;serializer:json;not null"`
	CategoryIDs  []uint     `gorm:"type:jsonb;serializer:json;not null"`
	Lat          float64    `gorm:"type:double precision;not null"`
	Long         float64    `gorm:"type:double precision;not null"`
	GeoHash      string     `gorm:"type:varchar(12);index"`
	Available    bool       `gorm:"default:true;not null"`
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}

// TokenCategoryFreq tracks how often a token appears across products in a given category.
// Owned exclusively by the search service. Populated by the background frequency flusher.
type TokenCategoryFreq struct {
	Token      string    `gorm:"primaryKey;type:varchar(200)"`
	CategoryID uint      `gorm:"primaryKey"`
	Count      int64     `gorm:"not null;default:0"`
	UpdatedAt  time.Time
}

