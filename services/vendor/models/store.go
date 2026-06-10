package models

import (
	"time"

	"gorm.io/gorm"
)

type Store struct {
	ID            uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string `gorm:"type:varchar(255);not null" json:"name"`
	StreetAddress string `gorm:"type:text;not null" json:"streetAddress"`

	BusinessID uint      `gorm:"index;not null" json:"businessId"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;" json:"business,omitempty"`

	Lat  float64 `gorm:"type:decimal(10,8);not null" json:"lat"`
	Long float64 `gorm:"type:decimal(11,8);not null" json:"long"`

	GeoHash string `gorm:"type:varchar(12);index:idx_stores_geohash" json:"geoHash"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewStore(name string, address string, lat float64, long float64) Store {
	return Store{
		Name:          name,
		StreetAddress: address,
		Lat:           lat,
		Long:          long,
	}
}
