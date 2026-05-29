package implementation

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BusinessRepository struct {
	db *gorm.DB
}

func NewBusinessRepository(db *gorm.DB) *BusinessRepository {
	return &BusinessRepository{db: db}
}

func (r *BusinessRepository) Create(ctx context.Context, business *models.Business) (*models.Business, error) {
	if err := r.db.WithContext(ctx).Create(business).Error; err != nil {
		return nil, err
	}
	return business, nil
}

func (r *BusinessRepository) GetByID(ctx context.Context, id uint) (*models.Business, error) {
	var business models.Business
	if err := r.db.WithContext(ctx).First(&business, id).Error; err != nil {
		return nil, err
	}
	return &business, nil
}

func (r *BusinessRepository) Update(ctx context.Context, businessID uint, businessName string, desc string) (*models.Business, error) {
	updates := map[string]interface{}{}

	if strings.TrimSpace(businessName) != "" {
		updates["business_name"] = businessName
	}

	if strings.TrimSpace(desc) != "" {
		updates["desc"] = desc
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	var business models.Business

	err := r.db.WithContext(ctx).
		Model(&models.Business{}).
		Clauses(clause.Returning{}).
		Where("id = ?", businessID).
		Updates(updates).
		Scan(&business).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update business: %w", err)
	}

	if business.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &business, nil
}

func (r *BusinessRepository) Delete(ctx context.Context, id uint) error {
	db := r.db.WithContext(ctx).Delete(&models.Business{}, id)
	if db.Error != nil {
		return fmt.Errorf("failed to delete business: %w", db.Error)
	}
	if db.RowsAffected == 0 {
		return fmt.Errorf("cannot delete: business record not found")
	}
	return nil
}

func (r *BusinessRepository) HardDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Delete(&models.Business{}, id).Error
}

func (r *BusinessRepository) GetByUserID(ctx context.Context, userID uint) (*models.Business, error) {
	var business models.Business
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&business).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("business not found for user ID %d", userID)
		}
		return nil, err
	}
	return &business, nil
}

func (r *BusinessRepository) GetByEmail(ctx context.Context, email string) (*models.Business, error) {
	var business models.Business
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&business).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("business not found with email %s", email)
		}
		return nil, err
	}
	return &business, nil
}

func (r *BusinessRepository) GetByPhone(ctx context.Context, phone string) (*models.Business, error) {
	var business models.Business
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&business).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("business not found with phone %s", phone)
		}
		return nil, err
	}
	return &business, nil
}

func (r *BusinessRepository) VerifyBusiness(ctx context.Context, id uint) error {
	db := r.db.WithContext(ctx).Model(&models.Business{}).
		Where("id = ?", id).
		Update("verified_business", true)

	if db.Error != nil {
		return fmt.Errorf("failed to verify business: %w", db.Error)
	}
	if db.RowsAffected == 0 {
		return fmt.Errorf("business not found with ID %d", id)
	}
	return nil
}

func (r *BusinessRepository) ToggleActiveStatus(ctx context.Context, id uint, isActive bool) error {
	db := r.db.WithContext(ctx).Model(&models.Business{}).
		Where("id = ?", id).
		Update("is_active", isActive)

	if db.Error != nil {
		return fmt.Errorf("failed to toggle business active status: %w", db.Error)
	}

	if db.RowsAffected == 0 {
		return fmt.Errorf("business not found with ID %d", id)
	}
	return nil
}
