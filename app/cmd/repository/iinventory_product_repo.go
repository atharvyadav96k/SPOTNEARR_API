package repository

import (
	"context"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
)

type IInventoryProduct interface {
	AddProduct(ctx context.Context, invProduct models.InventoryProduct, bizID uint) error
	GetInvProductList(ctx context.Context, invID uint, bizID uint) ([]models.Product, error)
	UpdateProduct(ctx context.Context, invProduct models.InventoryProduct, bizID uint) (models.InventoryProduct, error)
	RemoveProduct(ctx context.Context, invProdID uint, bizID uint) error
}
