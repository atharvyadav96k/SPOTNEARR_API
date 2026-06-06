package models

import "time"

// ProductToken tracks how many products in each category contain a given search token.
// Count is incremented when a product with that token is added to the category,
// and decremented when the product is deleted.
type ProductToken struct {
	Token      string    `gorm:"primaryKey;type:varchar(100)" json:"token"`
	CategoryID uint      `gorm:"primaryKey" json:"categoryId"`
	Category   *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
	Count      int       `gorm:"not null;default:1" json:"count"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
