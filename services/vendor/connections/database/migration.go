package database

import (
	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Business{},
		&models.Store{},
		&models.BusinessAccess{},
		&models.Category{},
		&models.Product{},
		&models.InventoryProduct{},
		&models.ProductToken{},
		&models.SearchSyncOutbox{},
	)
}
