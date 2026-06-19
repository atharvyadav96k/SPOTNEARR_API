package search

import "time"

// SearchEntry is the search index row. One row per inventory_product.
// Stores tokens and location for ranking, plus minimal display fields (name, price)
// so search results can be rendered without a second round-trip.
type SearchEntry struct {
	ID            uint       `gorm:"primaryKey"` // = inventory_product.id
	ProductID     uint       `gorm:"index;not null"`
	Name          string     `gorm:"type:varchar(255);not null;default:''"`
	Price         float64    `gorm:"type:decimal(12,2);not null;default:0"`
	PriceUnit     string     `gorm:"type:varchar(20);not null;default:''"`
	SearchTokens  []string   `gorm:"type:jsonb;serializer:json;not null"`
	CategoryIDs   []uint     `gorm:"type:jsonb;serializer:json;not null"`
	Lat           float64    `gorm:"type:double precision;not null"`
	Long          float64    `gorm:"type:double precision;not null"`
	GeoHash       string     `gorm:"type:varchar(12);index"`
	Available     bool       `gorm:"default:true;not null"`
	UpdatedAt     time.Time
	DeletedAt     *time.Time `gorm:"index"`
}

// TokenCategoryFreq tracks how often a token appears across products in a given category.
// Owned exclusively by the search service. Populated by the background frequency flusher.
type TokenCategoryFreq struct {
	Token      string    `gorm:"primaryKey;type:varchar(200)"`
	CategoryID uint      `gorm:"primaryKey"`
	Count      int64     `gorm:"not null;default:0"`
	UpdatedAt  time.Time
}

