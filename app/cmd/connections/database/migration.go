package database

import (
	userdb "github.com/Developer-Aadesh/spotnearr-database/user"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return userdb.AutoMigrate(db)
}
