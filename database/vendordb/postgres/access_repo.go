package postgres

import (
	"context"

	"github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"gorm.io/gorm"
)

type AccessRepository struct{ db *gorm.DB }

func NewAccessRepository(db *gorm.DB) *AccessRepository {
	return &AccessRepository{db: db}
}

func (a *AccessRepository) GetAccessByUserId(ctx context.Context, id uint) (vendordb.BusinessAccess, error) {
	var access vendordb.BusinessAccess
	db := a.db.WithContext(ctx).Where("user_id = ?", id).First(&access)
	if db.Error != nil {
		return access, db.Error
	}
	return access, nil
}

func (a *AccessRepository) GetAccessByBusinessId(ctx context.Context, id uint) (vendordb.BusinessAccess, error) {
	var access vendordb.BusinessAccess
	err := a.db.WithContext(ctx).Where("business_id = ?", id).First(&access).Error
	return access, err
}

func (a *AccessRepository) CreateNewAccess(ctx context.Context, access vendordb.BusinessAccess) error {
	return a.db.WithContext(ctx).Create(&access).Error
}

func (a *AccessRepository) RemoveAccessByUserId(ctx context.Context, id uint) error {
	return a.db.WithContext(ctx).Where("user_id = ?", id).Delete(&vendordb.BusinessAccess{}).Error
}
