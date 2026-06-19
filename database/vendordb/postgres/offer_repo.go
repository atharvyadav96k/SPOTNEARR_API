package postgres

import (
	"context"

	"github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OfferRepository struct{ db *gorm.DB }

func NewOfferRepository(db *gorm.DB) *OfferRepository {
	return &OfferRepository{db: db}
}

func (r *OfferRepository) Create(ctx context.Context, offer vendordb.Offer) (vendordb.Offer, error) {
	if err := r.db.WithContext(ctx).Create(&offer).Error; err != nil {
		return vendordb.Offer{}, err
	}
	return offer, nil
}

func (r *OfferRepository) GetByID(ctx context.Context, id uint, bizID uint) (vendordb.Offer, error) {
	var offer vendordb.Offer
	err := r.db.WithContext(ctx).
		Where("id = ? AND business_id = ?", id, bizID).
		First(&offer).Error
	return offer, err
}

func (r *OfferRepository) GetByBusiness(ctx context.Context, bizID uint) ([]vendordb.Offer, error) {
	var offers []vendordb.Offer
	err := r.db.WithContext(ctx).
		Where("business_id = ?", bizID).
		Order("created_at DESC").
		Find(&offers).Error
	return offers, err
}

func (r *OfferRepository) Update(ctx context.Context, offer vendordb.Offer, bizID uint) (vendordb.Offer, error) {
	updates := map[string]any{}
	if offer.Title != "" {
		updates["title"] = offer.Title
	}
	if offer.Description != "" {
		updates["description"] = offer.Description
	}
	if offer.DiscountType != "" {
		updates["discount_type"] = offer.DiscountType
	}
	if offer.DiscountValue > 0 {
		updates["discount_value"] = offer.DiscountValue
	}
	if offer.MinOrderValue != nil {
		updates["min_order_value"] = offer.MinOrderValue
	}
	if offer.ExpiresAt != nil {
		updates["expires_at"] = offer.ExpiresAt
	}
	if offer.Code != nil {
		updates["code"] = offer.Code
	}
	if offer.MaxUsage != nil {
		updates["max_usage"] = offer.MaxUsage
	}
	if len(updates) == 0 {
		return vendordb.Offer{}, nil
	}

	var updated vendordb.Offer
	db := r.db.WithContext(ctx).Model(&vendordb.Offer{}).
		Clauses(clause.Returning{}).
		Where("id = ? AND business_id = ?", offer.ID, bizID).
		Updates(updates).
		Scan(&updated)
	if db.Error != nil {
		return vendordb.Offer{}, db.Error
	}
	if db.RowsAffected == 0 {
		return vendordb.Offer{}, gorm.ErrRecordNotFound
	}
	return updated, nil
}

func (r *OfferRepository) Delete(ctx context.Context, id uint, bizID uint) error {
	db := r.db.WithContext(ctx).
		Where("id = ? AND business_id = ?", id, bizID).
		Delete(&vendordb.Offer{})
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
