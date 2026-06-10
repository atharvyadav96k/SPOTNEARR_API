package models

import (
	"time"

	"gorm.io/gorm"
)

type InventoryProduct struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	StoreID uint   `gorm:"uniqueIndex:idx_store_product;not null" json:"storeId"`
	Store   *Store `gorm:"foreignKey:StoreID;constraint:OnDelete:CASCADE;" json:"store,omitempty"`

	ProductID uint     `gorm:"uniqueIndex:idx_store_product;not null" json:"productId"`
	Product   *Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;" json:"product,omitempty"`

	Count *int `gorm:"default:null" json:"count"`

	Available *bool `gorm:"default:true" json:"available"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewInvProduct(storeID uint, productId uint, count *int, available *bool) InventoryProduct {
	return InventoryProduct{
		StoreID:   storeID,
		ProductID: productId,
		Count:     count,
		Available: available,
	}
}
