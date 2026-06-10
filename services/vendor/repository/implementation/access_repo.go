package implementation

import (
	"context"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
	"gorm.io/gorm"
)

type AccessRepository struct{ db *gorm.DB }

func NewAccessRepository(db *gorm.DB) *AccessRepository {
	return &AccessRepository{db: db}
}

func (a *AccessRepository) GetAccessByUserId(ctx context.Context, id uint) (models.BusinessAccess, error) {
	var access models.BusinessAccess
	db := a.db.WithContext(ctx).Where("user_id = ?", id).First(&access)
	if db.Error != nil {
		return access, db.Error
	}
	return access, nil
}

func (a *AccessRepository) GetAccessByBusinessId(ctx context.Context, id uint) (models.BusinessAccess, error) {
	var access models.BusinessAccess
	err := a.db.WithContext(ctx).Where("business_id = ?", id).First(&access).Error
	return access, err
}

func (a *AccessRepository) CreateNewAccess(ctx context.Context, access models.BusinessAccess) error {
	return a.db.WithContext(ctx).Create(&access).Error
}

func (a *AccessRepository) RemoveAccessByUserId(ctx context.Context, id uint) error {
	return a.db.WithContext(ctx).Where("user_id = ?", id).Delete(&models.BusinessAccess{}).Error
}
