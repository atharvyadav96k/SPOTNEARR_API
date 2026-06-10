package database

import (
	"log"

	dbmod "github.com/Developer-Aadesh/spotnearr-database"
	"gorm.io/gorm"
)

func InitDB(databaseUrl string) (*gorm.DB, error) {
	db, err := dbmod.Connect(databaseUrl)
	if err != nil {
		return nil, err
	}
	log.Println("Vendor Service: database connection established.")
	return db, nil
}
