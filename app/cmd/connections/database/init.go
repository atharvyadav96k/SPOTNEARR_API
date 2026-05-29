package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(databaseUrl string) (*gorm.DB, error) {
	DB, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		return nil, err
	}
	log.Println("Database connection successfully established.")
	return DB, nil
}
