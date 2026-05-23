package database

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Business{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.Store{}); err != nil {
		return err
	}
	return nil
}
