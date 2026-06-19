package postgres

import (
	"context"
	"math"

	"github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"gorm.io/gorm"
)

type SpotlightRepository struct{ db *gorm.DB }

func NewSpotlightRepository(db *gorm.DB) *SpotlightRepository {
	return &SpotlightRepository{db: db}
}

func (r *SpotlightRepository) Create(ctx context.Context, spotlight vendordb.Spotlight) (vendordb.Spotlight, error) {
	if err := r.db.WithContext(ctx).Create(&spotlight).Error; err != nil {
		return vendordb.Spotlight{}, err
	}
	return spotlight, nil
}

func (r *SpotlightRepository) GetByID(ctx context.Context, id uint, bizID uint) (vendordb.Spotlight, error) {
	var s vendordb.Spotlight
	err := r.db.WithContext(ctx).
		Where("id = ? AND business_id = ?", id, bizID).
		First(&s).Error
	return s, err
}

func (r *SpotlightRepository) GetByBusiness(ctx context.Context, bizID uint) ([]vendordb.Spotlight, error) {
	var spotlights []vendordb.Spotlight
	err := r.db.WithContext(ctx).
		Where("business_id = ?", bizID).
		Order("created_at DESC").
		Find(&spotlights).Error
	return spotlights, err
}

func (r *SpotlightRepository) Delete(ctx context.Context, id uint, bizID uint) error {
	db := r.db.WithContext(ctx).
		Where("id = ? AND business_id = ?", id, bizID).
		Delete(&vendordb.Spotlight{})
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SpotlightRepository) GetFeed(ctx context.Context, lat, long, rangeKm float64) ([]vendordb.Spotlight, error) {
	latDelta := rangeKm / 111.32
	lonDelta := rangeKm / (111.32 * math.Cos(lat*math.Pi/180))
	minLat, maxLat := lat-latDelta, lat+latDelta
	minLon, maxLon := long-lonDelta, long+lonDelta

	var spotlights []vendordb.Spotlight
	err := r.db.WithContext(ctx).
		Distinct("spotlights.*").
		Joins("JOIN stores s ON s.business_id = spotlights.business_id AND s.deleted_at IS NULL").
		Where("spotlights.deleted_at IS NULL").
		Where("s.lat BETWEEN ? AND ? AND s.long BETWEEN ? AND ?", minLat, maxLat, minLon, maxLon).
		Order("spotlights.created_at DESC").
		Preload("Business").
		Limit(50).
		Find(&spotlights).Error
	return spotlights, err
}
