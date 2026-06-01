package database

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Business{},
		&models.Store{},
		&models.BusinessAccess{},
		&models.Product{},
		&models.InventoryProduct{},
		&models.Claim{},
	); err != nil {
		return err
	}
	return nil
}
