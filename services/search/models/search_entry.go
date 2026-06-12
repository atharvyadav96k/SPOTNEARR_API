package models

import "time"

// SearchEntry mirrors the schema in the database module.
// One row per inventory_product; stores tokens, location, and minimal display fields.
type SearchEntry struct {
	ID            uint       `gorm:"primaryKey"`
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
