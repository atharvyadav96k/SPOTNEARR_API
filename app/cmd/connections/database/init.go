package database

import (
	dbmod "github.com/Developer-Aadesh/spotnearr-database"
	"gorm.io/gorm"
)

func InitDB(databaseURL string) (*gorm.DB, error) {
	return dbmod.Connect(databaseURL)
}
