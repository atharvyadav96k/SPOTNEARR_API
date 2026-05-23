package app

import "gorm.io/gorm"

func (a *App) GetDb() *gorm.DB {
	return a.db
}
