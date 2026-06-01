package models

import (
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

	InventoryProductID uint              `gorm:"index;not null" json:"inventoryProductId"`
	InventoryProduct   *InventoryProduct `gorm:"foreignKey:InventoryProductID;constraint:OnDelete:CASCADE;" json:"inventoryProduct,omitempty"`

	Status ClaimStatus `gorm:"type:varchar(20);default:'pending';not null" json:"status"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewClaim(userID uint, invProductID uint) Claim {
	return Claim{
		UserID:             userID,
		InventoryProductID: invProductID,
		Status:             ClaimStatusPending,
	}
}
