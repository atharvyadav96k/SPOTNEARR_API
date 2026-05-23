package app

import "gorm.io/gorm"

type envs struct {
	DATABASE_URL string
}

type App struct {
	db *gorm.DB
}
