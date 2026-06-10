package database

import (
	searchdb "github.com/Developer-Aadesh/spotnearr-database/search"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return searchdb.AutoMigrate(db)
}
