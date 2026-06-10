package models

import (
	"time"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/utils/request"
	"gorm.io/gorm"
)

type Product struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	Name      string  `gorm:"type:varchar(255);not null" json:"name"`
	Price     float64 `gorm:"type:decimal(12,2);not null" json:"price"`
	PriceUnit string  `gorm:"type:varchar(20);not null" json:"priceUnit"`

	Quantity     *float64 `gorm:"type:decimal(12,2)" json:"quantity"`
	QuantityUnit *string  `gorm:"type:varchar(20)" json:"quantityUnit"`
	Desc         string   `gorm:"type:text" json:"desc"`

	BusinessID uint      `gorm:"index;not null" json:"businessId"`
	Business   *Business `gorm:"foreignKey:BusinessID;constraint:OnDelete:CASCADE;" json:"business,omitempty"`

	SearchTokens []string   `gorm:"type:jsonb;serializer:json" json:"-"`
	Categories   []Category `gorm:"many2many:product_categories;" json:"categories,omitempty"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`
}

func NewProduct(name string, price request.ValueUnit, desc string, quantity request.ValueUnit, searchTokens []string) Product {
	return Product{
		Name:         name,
		Price:        price.Value,
		PriceUnit:    price.Unit,
		Desc:         desc,
		Quantity:     &quantity.Value,
		QuantityUnit: &quantity.Unit,
		SearchTokens: searchTokens,
	}
}
