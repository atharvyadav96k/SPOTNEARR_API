package models

import (
	"time"

	"gorm.io/gorm"
)

type Business struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	BusinessName string `gorm:"type:varchar(255);not null"`

	UserID uint  `gorm:"index;unique;not null"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	Email *string `gorm:"type:varchar(255);unique;index"`
	Phone *string `gorm:"type:varchar(20);unique;index"`

	IsActive         bool `gorm:"default:true;not null"`
	VerifiedBusiness bool `gorm:"default:false;not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
