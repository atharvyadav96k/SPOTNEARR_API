package database

import (
	"log"

	dbmod "github.com/Developer-Aadesh/spotnearr-database"
	"gorm.io/gorm"
)

func InitDB(url string) (*gorm.DB, error) {
	db, err := dbmod.Connect(url)
	if err != nil {
		return nil, err
	}
	log.Println("Search Service: database connection established.")
	return db, nil
}
