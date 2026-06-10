package models

import "time"

type ProductToken struct {
	Token      string    `gorm:"primaryKey;type:varchar(100)" json:"token"`
	CategoryID uint      `gorm:"primaryKey" json:"categoryId"`
	Category   *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
	Count      int       `gorm:"not null;default:1" json:"count"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
