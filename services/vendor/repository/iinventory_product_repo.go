package repository

import (
	"context"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
)

type IInventoryProduct interface {
	AddProduct(ctx context.Context, invProduct models.InventoryProduct, bizID uint) error
	GetInvProductList(ctx context.Context, invID uint, bizID uint) ([]models.Product, error)
	UpdateProduct(ctx context.Context, invProduct models.InventoryProduct, bizID uint) (models.InventoryProduct, error)
	RemoveProduct(ctx context.Context, invProdID uint, bizID uint) error
	GetByID(ctx context.Context, id uint) (*InvProductDetail, error)
}

// InvProductDetail is the snapshot returned by the internal endpoint for claim validation.
type InvProductDetail struct {
	ID            uint    `json:"id"`
	Available     bool    `json:"available"`
	ProductName   string  `json:"product_name"`
	Price         float64 `json:"price"`
	PriceUnit     string  `json:"price_unit"`
	StoreName     string  `json:"store_name"`
	StreetAddress string  `json:"street_address"`
	Lat           float64 `json:"lat"`
	Long          float64 `json:"long"`
}
