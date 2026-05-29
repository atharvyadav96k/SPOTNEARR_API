package app

import (
	"os"

	"github.com/atharvyadav96k/SPOTNEARR_API/database"
)

func (a *App) InitDb() error {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// local database (development)
		dbURL = "postgresql://admin:admin123@localhost:5432/spotnearr"
	}
	a.db, err = database.InitDB(dbURL)
	if err != nil {
		return err
	}
	if err := database.AutoMigrate(a.db); err != nil {
		return err
	}
	return nil
}
