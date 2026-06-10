package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(databaseUrl string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	log.Println("Vendor Service: database connection established.")
	return db, nil
}
