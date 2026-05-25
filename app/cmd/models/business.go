package models

import (
	"time"

	"gorm.io/gorm"
)

type Business struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	BusinessName string `gorm:"type:varchar(255);not null" json:"businessName"`

	UserID uint  `gorm:"index;unique;not null" json:"userId"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`

	Email *string `gorm:"type:varchar(255);unique;index" json:"email"`
	Phone *string `gorm:"type:varchar(20);unique;index" json:"phone"`

	IsActive         bool `gorm:"default:true;not null" json:"isActive"`
	VerifiedBusiness bool `gorm:"default:false;not null" json:"verifiedBusiness"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
