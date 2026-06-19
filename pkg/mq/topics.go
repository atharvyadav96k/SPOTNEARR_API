package mq

import "time"

// Topic is a routing key on the spotnearr.events exchange.
type Topic string

const (
	TopicBusinessFollow   Topic = "user.business.follow"
	TopicBusinessUnfollow Topic = "user.business.unfollow"
	TopicProductSync      Topic = "vendor.product.sync"
	TopicDealSync         Topic = "vendor.deal.sync"
	TopicDealDelete       Topic = "vendor.deal.delete"
)

// BusinessFollowPayload is published for both follow and unfollow events.
type BusinessFollowPayload struct {
	BusinessID uint `json:"business_id"`
	UserID     uint `json:"user_id"`
}

// DealStoreLocation carries the geohash5 cell for one store location.
type DealStoreLocation struct {
	StoreID  uint   `json:"store_id"`
	GeoHash5 string `json:"geohash5"`
}

// DealSyncPayload is published when a vendor creates or updates an offer.
type DealSyncPayload struct {
	DealID        uint                `json:"deal_id"`
	BusinessID    uint                `json:"business_id"`
	Title         string              `json:"title"`
	DiscountType  string              `json:"discount_type"`
	DiscountValue float64             `json:"discount_value"`
	MinOrderValue *float64            `json:"min_order_value,omitempty"`
	Active        bool                `json:"active"`
	ExpiresAt     *time.Time          `json:"expires_at,omitempty"`
	Locations     []DealStoreLocation `json:"locations"`
}

// DealDeletePayload is published when a vendor deletes an offer.
type DealDeletePayload struct {
	DealID uint `json:"deal_id"`
}
