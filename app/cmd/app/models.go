package app

import (
	"github.com/atharvyadav96k/SPOTNEARR_API/connections/cache"
	"gorm.io/gorm"
)

type envs struct {
	DATABASE_URL string
}

type App struct {
	db    *gorm.DB
	cache *cache.Cache
}
