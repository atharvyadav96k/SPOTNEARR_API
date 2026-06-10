package database

import (
	vendordb "github.com/Developer-Aadesh/spotnearr-database/vendordb"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := vendordb.AutoMigrate(db); err != nil {
		return err
	}
	return vendordb.RunIndexes(db)
}
