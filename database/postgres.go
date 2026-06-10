package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect opens a GORM postgres connection. All services use this instead of
// duplicating the same gorm.Open boilerplate.
func Connect(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	log.Println("database connection established")
	return db, nil
}
