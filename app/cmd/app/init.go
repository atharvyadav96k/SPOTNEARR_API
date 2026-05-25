package app

import "github.com/atharvyadav96k/SPOTNEARR_API/database"

func (a *App) InitDb() error {
	var err error
	a.db, err = database.InitDB("DATABASE_URL")
	if err != nil {
		return err
	}
	if err := database.AutoMigrate(a.db); err != nil {
		return err
	}
	return nil
}
