package models

import "time"

// SearchEntry mirrors the lean schema in the database module.
// One row per inventory_product; only search-critical fields are stored.
type SearchEntry struct {
	ID           uint       `gorm:"primaryKey"`
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
