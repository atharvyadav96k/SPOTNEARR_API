package models

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"type:varchar(255);not null"`
	StreetAddress string `gorm:"type:text;not null"`

	BusinessID uint      `gorm:"index;not null"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;"`

	Lat  float64 `gorm:"type:decimal(10,8);not null"`
	Long float64 `gorm:"type:decimal(11,8);not null"`

	GeoHash string `gorm:"type:varchar(12);index:idx_stores_geohash"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
