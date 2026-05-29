package app

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"gorm.io/gorm"
)

func (a *App) GetDb() *gorm.DB {
	return a.db
}

func (a *App) GetCache() *cache.Cache {
	return a.cache
}
