package implementation

import (
	"context"
	"errors"
	"fmt"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (s *StoreRepository) GetStoreByBusinessId(ctx context.Context, id uint) ([]models.Store, error) {
	type joinedResult struct {
		models.Store
		BusinessIDRef *uint `gorm:"column:biz_id"`
	}
	var results []joinedResult
	db := s.db.WithContext(ctx).
		Table("businesses").
		Select("stores.*, businesses.id as biz_id").
		Joins("LEFT JOIN stores ON stores.business_id = businesses.id AND stores.deleted_at IS NULL").
		Where("businesses.id = ? AND businesses.deleted_at IS NULL", id).
		Scan(&results)

	if db.Error != nil {
		return nil, db.Error
	}
	if db.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var stores []models.Store
	for _, res := range results {
		if res.Store.ID == 0 {
			continue
		}
		stores = append(stores, res.Store)
	}
	return stores, nil
}

func (s *StoreRepository) CreateStoreByBusinessId(ctx context.Context, store *models.Store) error {
	if err := s.db.WithContext(ctx).Create(store).Error; err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}
	return nil
}

func (s *StoreRepository) UpdateStoreByBusinessId(ctx context.Context, store models.Store) (models.Store, error) {
	var existingStore models.Store
	err := s.db.WithContext(ctx).
		Where("id = ? AND business_id = ?", store.ID, store.BusinessID).
		First(&existingStore).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Store{}, fmt.Errorf("store with ID %d not found for business %d", store.ID, store.BusinessID)
		}
		return models.Store{}, err
	}

	if err := s.db.WithContext(ctx).Save(store).Error; err != nil {
		return models.Store{}, fmt.Errorf("failed to update store values: %w", err)
	}

	return store, nil
}
