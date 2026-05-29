package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	Name  string  `gorm:"type:varchar(255);not null" json:"name"`
	Price float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Desc  string  `gorm:"type:text" json:"desc"`

	BusinessID uint      `gorm:"index;not null" json:"businessId"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;" json:"business,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewProduct(name string, price float64) Product {
	return Product{
		Name:  name,
		Price: price,
	}
}
