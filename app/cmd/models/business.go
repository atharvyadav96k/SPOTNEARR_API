package models

import (
	"time"

	"gorm.io/gorm"
)

type Business struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	BusinessName string `gorm:"type:varchar(255);not null" json:"businessName"`

	Email *string `gorm:"type:varchar(255);unique;index" json:"email"`
	Phone *string `gorm:"type:varchar(20);unique;index" json:"phone"`
	Desc  string  `gorm:"type:varchar(100)" json:"description"`

	IsActive         bool `gorm:"default:true;not null" json:"isActive"`
	VerifiedBusiness bool `gorm:"default:false;not null" json:"verifiedBusiness"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewBusiness(name string, email string, phone string, desc string) *Business {
	return &Business{
		BusinessName: name,
		Email:        &email,
		Phone:        &phone,
		Desc:         desc,
	}
}
