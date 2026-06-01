package app

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/config"
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/database"
)

func (a *App) InitDb() error {
	var err error
	a.db, err = database.InitDB(config.C.DatabaseURL)
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
	a.cache, err = cache.InitCache(config.C.CacheURL, config.C.CachePassword)
	if err != nil {
		return err
	}
	return nil
}

func (a *App) InitCaptcha() error {
	return nil
}
