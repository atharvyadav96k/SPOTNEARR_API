package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type InventoryService struct {
	base_service
}

func NewInventoryService(db *gorm.DB, cache *cache.Cache) *InventoryService {
	return &InventoryService{
		base_service: NewBaseService(db, cache),
	}
}

func (i *InventoryService) GetBusinessInventory(id uint) response.Res {
	stores, err := i.RepoStore().GetStoreByBusinessId(context.Background(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i.ResponseNotFound("Business not found")
		}
		return i.ResponseBadRequest(err.Error())
	}
	return i.ResponseOK("Business inventories", stores)
}

func (i *InventoryService) CreateInventory(id uint, store *models.Store) response.Res {
	store.BusinessID = id
	err := i.RepoStore().CreateStoreByBusinessId(context.Background(), store)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return i.ResponseConflict("Failed to create store due to invalid store")
		}
		return i.ResponseBadRequest("Failed to create store")
	}
	return i.ResponseOK("Store registered successfully", nil)
}

func (i *InventoryService) UpdateInventory(businessId uint, storeId uint, name string, address string, lat float64, long float64) response.Res {
	store := models.NewStore(name, address, lat, long)
	store, err := i.RepoStore().UpdateStoreByBusinessId(context.Background(), store)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i.ResponseNotFound("Store not found")
		}
		return i.ResponseInternalServer("Failed to update store")
	}
	return i.ResponseOK("store updated successfully", store)
}

func (i *InventoryService) GetInvProducts(bizID uint, storeID uint) response.Res {
	products, err := i.RepoInvProduct().GetInvProductList(context.Background(), storeID, bizID)
	if err != nil {
		return i.ResponseInternalServer("Failed to load the product")
	}
	return i.ResponseOK("Products list", products)
}

func (i *InventoryService) AddInvProduct(bizID uint, invID uint, productID uint, count *int, available *bool) response.Res {
	product := models.NewInvProduct(invID, productID, count, available)
	err := i.RepoInvProduct().AddProduct(context.Background(), product, bizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i.ResponseNotFound("Product not found.")
		}
		return i.ResponseInternalServer("Failed to add the product in store.")
	}
	return i.ResponseOK("Product added successfully to store", nil)
}
