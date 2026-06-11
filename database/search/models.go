package search

import "time"

// SearchEntry is the flat denormalized read table. One row per inventory_product.
// Written by the sync worker (outbox consumer); never written by the search query path.
type SearchEntry struct {
	ID            uint       `gorm:"primaryKey" json:"id"` // = inventory_product.id
	ProductID     uint       `gorm:"index;not null" json:"productId"`
	BusinessID    uint       `gorm:"index;not null" json:"businessId"`
	ProductName   string     `gorm:"type:text;not null" json:"productName"`
	Price         float64    `gorm:"type:numeric(12,2);not null" json:"price"`
	PriceUnit     string     `gorm:"type:varchar(20);not null" json:"priceUnit"`
	Quantity      *float64   `gorm:"type:numeric(12,2)" json:"quantity,omitempty"`
	QuantityUnit  *string    `gorm:"type:varchar(20)" json:"quantityUnit,omitempty"`
	Description   string     `gorm:"type:text" json:"description"`
	SearchTokens  []string   `gorm:"type:jsonb;serializer:json;not null" json:"-"`
	CategoryIDs   []uint     `gorm:"type:jsonb;serializer:json;not null" json:"-"`
	StoreID       uint       `gorm:"index;not null" json:"storeId"`
	StoreName     string     `gorm:"type:text;not null" json:"storeName"`
	StreetAddress string     `gorm:"type:text;not null" json:"streetAddress"`
	Lat           float64    `gorm:"type:double precision;not null" json:"lat"`
	Long          float64    `gorm:"type:double precision;not null" json:"long"`
	GeoHash       string     `gorm:"type:varchar(12);index" json:"geoHash"`
	Available     bool       `gorm:"default:true;not null" json:"available"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `gorm:"index" json:"deletedAt,omitempty"`
}

// ProductResult is the shape returned to API callers from a search query.
type ProductResult struct {
	ID            uint     `json:"id"`
	ProductID     uint     `json:"productId"`
	BusinessID    uint     `json:"businessId"`
	ProductName   string   `json:"productName"`
	Price         float64  `json:"price"`
	PriceUnit     string   `json:"priceUnit"`
	Quantity      *float64 `json:"quantity,omitempty"`
	QuantityUnit  *string  `json:"quantityUnit,omitempty"`
	Description   string   `json:"description"`
	StoreID       uint     `json:"storeId"`
	StoreName     string   `json:"storeName"`
	StreetAddress string   `json:"streetAddress"`
	Lat           float64  `json:"lat"`
	Long          float64  `json:"long"`
	GeoHash       string   `json:"geoHash"`
	DistanceKm    float64  `json:"distanceKm,omitempty"`
	TokenMatchCnt int      `json:"-"`
	CategoryIDs   []uint   `json:"-"`
}
