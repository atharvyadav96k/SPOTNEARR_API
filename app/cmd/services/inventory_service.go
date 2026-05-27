package services

import (
	"context"
	"errors"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/repository"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"gorm.io/gorm"
)

type InventoryService struct {
	base_service
	repo repository.IStoreRepository
}

func NewInventoryService(db *gorm.DB) *InventoryService {
	return &InventoryService{
		base_service: NewBaseService(db),
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
	store := models.Store{
		Name:          name,
		StreetAddress: address,
		Lat:           lat,
		Long:          long,
	}
	store, err := i.RepoStore().UpdateStoreByBusinessId(context.Background(), store)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return i.ResponseNotFound("Store not found")
		}
		return i.ResponseInternalServer("Failed to update store")
	}
	return i.ResponseOK("store updated successfully", store)
}
