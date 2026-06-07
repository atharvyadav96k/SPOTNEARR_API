package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type ClaimStatus string

const (
	ClaimStatusPending  ClaimStatus = "pending"
	ClaimStatusAccepted ClaimStatus = "accepted"
	ClaimStatusRejected ClaimStatus = "rejected"
)

type Claim struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint  `gorm:"index;not null" json:"userId"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`

	InventoryProductID uint `gorm:"index;not null" json:"inventoryProductId"`

	Status ClaimStatus `gorm:"type:varchar(20);default:'pending';not null" json:"status"`

	// Frozen snapshot of the inventory product at claim time.
	// Populated via Vendor Service HTTP call so GET /claims never needs a cross-service call.
	InvProductSnapshot json.RawMessage `gorm:"type:jsonb" json:"product,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewClaim(userID uint, invProductID uint, snapshot json.RawMessage) Claim {
	return Claim{
		UserID:             userID,
		InventoryProductID: invProductID,
		Status:             ClaimStatusPending,
		InvProductSnapshot: snapshot,
	}
}
