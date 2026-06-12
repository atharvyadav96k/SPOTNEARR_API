package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	vendormodel "github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"gorm.io/gorm"
)

type ProductService struct{ baseService }

func NewProductService(db *gorm.DB, c *cache.Cache) *ProductService {
	return &ProductService{baseService: newBaseService(db, c)}
}

func (p *ProductService) GetProductByID(bizID uint, productID uint) httputil.Res {
	product, err := p.RepoProduct().GetProductById(context.Background(), productID, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.ResponseNotFound("Product not found")
		}
		return p.ResponseInternalServer("Failed to get product")
	}
	return p.ResponseOK("Product", product)
}

func (p *ProductService) GetAllProducts(bizID uint) httputil.Res {
	products, err := p.RepoProduct().GetProductsByBusiness(context.Background(), bizID)
	if err != nil {
		return p.ResponseInternalServer("Failed to list products")
	}
	return p.ResponseOK("Products", products)
}

func (p *ProductService) AddNewProduct(bizID uint, product vendormodel.Product, storeIDs []uint) httputil.Res {
	ctx := context.Background()
	created, err := p.RepoProduct().AddProduct(ctx, product, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return p.ResponseBadRequest("Business does not exist")
		}
		return p.ResponseInternalServer("Failed to add product")
	}

	// Auto-link to first store if no storeIDs provided.
	if len(storeIDs) == 0 {
		stores, storeErr := p.RepoStore().GetStoreByBusinessId(ctx, bizID)
		if storeErr == nil && len(stores) > 0 {
			storeIDs = []uint{stores[0].ID}
		}
	}

	available := true
	for _, storeID := range storeIDs {
		invProduct := vendormodel.NewInvProduct(storeID, created.ID, nil, &available)
		_ = p.RepoInvProduct().AddProduct(ctx, invProduct, bizID)
	}

	return p.ResponseCreated("Product added successfully", created)
}

func (p *ProductService) UpdateProduct(bizID uint, product vendormodel.Product) httputil.Res {
	ctx := context.Background()
	updated, err := p.RepoProduct().UpdateProduct(ctx, product, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.ResponseNotFound("Product not found")
		}
		return p.ResponseInternalServer("Failed to update product")
	}
	// Write outbox entries so the search index reflects the new name/price.
	_ = p.RepoInvProduct().WriteUpsertOutboxesForProduct(ctx, updated.ID)
	return p.ResponseOK("Product updated successfully", updated)
}

func (p *ProductService) DeleteProduct(bizID uint, productID uint) httputil.Res {
	ctx := context.Background()
	// Write delete outbox entries before the soft-delete so inv_product IDs are still queryable.
	_ = p.RepoInvProduct().WriteDeleteOutboxesForProduct(ctx, productID)
	if err := p.RepoProduct().DeleteProduct(ctx, productID, bizID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return p.ResponseNotFound("Product not found")
		}
		return p.ResponseInternalServer("Failed to delete product")
	}
	return p.ResponseOK("Product removed successfully", nil)
}
