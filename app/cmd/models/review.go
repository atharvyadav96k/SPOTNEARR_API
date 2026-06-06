package models

import (
	"time"

	"gorm.io/gorm"
)

type ReviewTarget string

const (
	ReviewTargetBusiness  ReviewTarget = "business"
	ReviewTargetProduct   ReviewTarget = "product"
	ReviewTargetOffer     ReviewTarget = "offer"
	ReviewTargetSpotlight ReviewTarget = "spotlight"
)

type Review struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	UserID uint  `gorm:"index;not null" json:"userId"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`

	TargetType ReviewTarget `gorm:"type:varchar(20);not null;index:idx_review_target" json:"targetType"`
	TargetID   uint         `gorm:"not null;index:idx_review_target" json:"targetId"`

	Stars   uint8  `gorm:"not null" json:"stars"`
	Comment string `gorm:"type:text" json:"comment"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}
