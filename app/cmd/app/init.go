package app

import (
	"log"
	"os"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/database"
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

func (a *App) InitCache() error {
	var err error
	cacheUrl := os.Getenv("CACHE_URL")
	password := os.Getenv("CACHE_PASSWORD")
	if strings.TrimSpace(cacheUrl) == "" {
		cacheUrl = "localhost:6380"
	}
	log.Default().Println(cacheUrl, password)
	a.cache, err = cache.InitCache(cacheUrl, password)
	if err != nil {
		return err
	}
	return nil
}
