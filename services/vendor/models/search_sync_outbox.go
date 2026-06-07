package models

import "time"

// SearchSyncOutbox is the transactional outbox for propagating product/inventory
// changes to the Search Service. Each row is written in the same DB transaction
// as the triggering write (add/delete inventory product). The Search Service
// polls this table and applies events idempotently to search_entries.
type SearchSyncOutbox struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	EventType   string     `gorm:"type:varchar(20);not null" json:"eventType"` // "upsert" | "delete"
	Payload     []byte     `gorm:"type:jsonb;not null" json:"payload"`
	Processed   bool       `gorm:"default:false;index" json:"processed"`
	ProcessedAt *time.Time `json:"processedAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// OutboxPayload is the JSON structure stored in SearchSyncOutbox.Payload.
type OutboxPayload struct {
	EventType     string    `json:"event_type"`
	InvProductID  uint      `json:"inv_product_id"`
	ProductID     uint      `json:"product_id,omitempty"`
	BusinessID    uint      `json:"business_id,omitempty"`
	ProductName   string    `json:"product_name,omitempty"`
	Price         float64   `json:"price,omitempty"`
	PriceUnit     string    `json:"price_unit,omitempty"`
	Quantity      *float64  `json:"quantity,omitempty"`
	QuantityUnit  *string   `json:"quantity_unit,omitempty"`
	Desc          string    `json:"desc,omitempty"`
	SearchTokens  []string  `json:"search_tokens,omitempty"`
	Categories    []CatRef  `json:"categories,omitempty"`
	StoreID       uint      `json:"store_id,omitempty"`
	StoreName     string    `json:"store_name,omitempty"`
	StreetAddress string    `json:"street_address,omitempty"`
	Lat           float64   `json:"lat,omitempty"`
	Long          float64   `json:"long,omitempty"`
	GeoHash       string    `json:"geo_hash,omitempty"`
	Available     bool      `json:"available"`
}

type CatRef struct {
	ID uint `json:"id"`
}
