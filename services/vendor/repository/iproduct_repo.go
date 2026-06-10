package repository

import (
	"context"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
)

type IProductRepository interface {
	AddProduct(ctx context.Context, product models.Product, bizID uint) (models.Product, error)
	UpdateProduct(ctx context.Context, product models.Product, bizID uint) (models.Product, error)
	GetProductsByBusiness(ctx context.Context, bizID uint) ([]models.Product, error)
	GetProductById(ctx context.Context, id uint, bizID uint) (models.Product, error)
	DeleteProduct(ctx context.Context, productID uint, bizID uint) error
}
