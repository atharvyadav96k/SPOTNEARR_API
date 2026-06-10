package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/connections/cache"
	vendormodel "github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"gorm.io/gorm"
)

type InventoryService struct{ baseService }

func NewInventoryService(db *gorm.DB, c *cache.Cache) *InventoryService {
	return &InventoryService{baseService: newBaseService(db, c)}
}

func (i *InventoryService) GetBusinessInventory(bizID uint) httputil.Res {
	stores, err := i.RepoStore().GetStoreByBusinessId(context.Background(), bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i.ResponseNotFound("Business not found")
		}
		return i.ResponseBadRequest(err.Error())
	}
	return i.ResponseOK("Business inventories", stores)
}

func (i *InventoryService) CreateInventory(bizID uint, store *vendormodel.Store) httputil.Res {
	store.BusinessID = bizID
	if err := i.RepoStore().CreateStoreByBusinessId(context.Background(), store); err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return i.ResponseConflict("Failed to create store due to invalid business")
		}
		return i.ResponseBadRequest("Failed to create store")
	}
	return i.ResponseOK("Store registered successfully", nil)
}

func (i *InventoryService) UpdateInventory(bizID uint, storeID uint, name string, address string, lat float64, long float64) httputil.Res {
	store := vendormodel.NewStore(name, address, lat, long)
	store.ID = storeID
	store.BusinessID = bizID
	updated, err := i.RepoStore().UpdateStoreByBusinessId(context.Background(), store)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i.ResponseNotFound("Store not found")
		}
		return i.ResponseInternalServer("Failed to update store")
	}
	return i.ResponseOK("Store updated successfully", updated)
}

func (i *InventoryService) GetInvProducts(bizID uint, storeID uint) httputil.Res {
	products, err := i.RepoInvProduct().GetInvProductList(context.Background(), storeID, bizID)
	if err != nil {
		return i.ResponseInternalServer("Failed to load products")
	}
	return i.ResponseOK("Products list", products)
}

func (i *InventoryService) AddInvProduct(bizID uint, storeID uint, productID uint, count *int, available *bool) httputil.Res {
	product := vendormodel.NewInvProduct(storeID, productID, count, available)
	if err := i.RepoInvProduct().AddProduct(context.Background(), product, bizID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i.ResponseNotFound("Store not found")
		}
		return i.ResponseInternalServer("Failed to add product to store")
	}
	return i.ResponseOK("Product added successfully to store", nil)
}

func (i *InventoryService) UpdateInvProduct(bizID uint, invProdID uint, count *int, available *bool) httputil.Res {
	invProd := vendormodel.InventoryProduct{
		ID:        invProdID,
		Count:     count,
		Available: available,
	}
	updated, err := i.RepoInvProduct().UpdateProduct(context.Background(), invProd, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i.ResponseNotFound("Inventory product not found")
		}
		return i.ResponseInternalServer("Failed to update inventory product")
	}
	return i.ResponseOK("Inventory product updated", updated)
}

func (i *InventoryService) RemoveInvProduct(bizID uint, invProdID uint) httputil.Res {
	if err := i.RepoInvProduct().RemoveProduct(context.Background(), invProdID, bizID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i.ResponseNotFound("Inventory product not found")
		}
		return i.ResponseInternalServer("Failed to remove inventory product")
	}
	return i.ResponseOK("Product removed from store", nil)
}
