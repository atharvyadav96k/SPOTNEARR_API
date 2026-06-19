package postgres

import (
	"context"

	userdb "github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DealRepository struct{ db *gorm.DB }

func NewDealRepository(db *gorm.DB) *DealRepository {
	return &DealRepository{db: db}
}

func (r *DealRepository) Upsert(ctx context.Context, entry userdb.DealEntry) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "deal_id"}, {Name: "geo_hash5"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name", "price", "deal_price", "discount_type", "active", "expires_at", "updated_at",
			}),
		}).
		Create(&entry).Error
}

func (r *DealRepository) DeleteByDealID(ctx context.Context, dealID uint) error {
	return r.db.WithContext(ctx).
		Where("deal_id = ?", dealID).
		Delete(&userdb.DealEntry{}).Error
}

func (r *DealRepository) GetByIDs(ctx context.Context, dealIDs []uint) ([]userdb.DealEntry, error) {
	var entries []userdb.DealEntry
	err := r.db.WithContext(ctx).
		Where("deal_id IN ? AND active = true", dealIDs).
		Where("expires_at IS NULL OR expires_at > NOW()").
		Find(&entries).Error
	return entries, err
}

func (r *DealRepository) GetByGeoHash5(ctx context.Context, geohash5 string) ([]userdb.DealEntry, error) {
	var entries []userdb.DealEntry
	err := r.db.WithContext(ctx).
		Where("geo_hash5 = ? AND active = true", geohash5).
		Where("expires_at IS NULL OR expires_at > NOW()").
		Find(&entries).Error
	return entries, err
}

func (r *DealRepository) GetGeoHash5ForDeal(ctx context.Context, dealID uint) ([]string, error) {
	var hashes []string
	err := r.db.WithContext(ctx).
		Model(&userdb.DealEntry{}).
		Where("deal_id = ?", dealID).
		Pluck("geo_hash5", &hashes).Error
	return hashes, err
}
