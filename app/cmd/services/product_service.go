package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type ProductService struct {
	base_service
}

func NewProductService(db *gorm.DB, cache *cache.Cache) *ProductService {
	return &ProductService{
		base_service: NewBaseService(db, cache),
	}
}

func (p ProductService) GetProductByID(bizID uint, productID uint) response.Res {
	product, err := p.RepoProduct().GetProductById(context.Background(), productID, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.ResponseNotFound("product not found")
		}
		return p.ResponseInternalServer("Failed to get product")
	}
	return p.ResponseOK("Product", product)
}

func (p *ProductService) AddNewProduct(bizID uint, product models.Product, storeIDs []uint) response.Res {
	ctx := context.Background()
	created, err := p.RepoProduct().AddProduct(ctx, product, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return p.ResponseBadRequest("Business does not exist")
		}
		return p.ResponseInternalServer("Failed to add product")
	}

	if len(storeIDs) == 0 {
		stores, storeErr := p.RepoStore().GetStoreByBusinessId(ctx, bizID)
		if storeErr == nil && len(stores) > 0 {
			storeIDs = []uint{stores[0].ID}
		}
	}

	available := true
	for _, storeID := range storeIDs {
		invProduct := models.NewInvProduct(storeID, created.ID, nil, &available)
		_ = p.RepoInvProduct().AddProduct(ctx, invProduct, bizID)
	}

	return p.ResponseCreated("Product added successfully", created)
}

func (p *ProductService) UpdateProduct(bizID uint, product models.Product) response.Res {
	product, err := p.RepoProduct().UpdateProduct(context.Background(), product, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.ResponseBadRequest("Product not found")
		}
		return p.ResponseInternalServer("Failed to update product")
	}
	return p.ResponseOK("Product updated successfully", product)
}

func (p *ProductService) DeleteProduct(bizID uint, productID uint) response.Res {
	err := p.RepoProduct().DeleteProduct(context.Background(), productID, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.ResponseNotFound("Product not found")
		}
		return p.ResponseInternalServer("Failed to delete product")
	}
	return p.ResponseOK("Product remove successfully", nil)
}
