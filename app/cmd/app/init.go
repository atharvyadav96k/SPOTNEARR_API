package app

import "github.com/atharvyadav96k/SPOTNEARR_API/database"

func (a *App) InitDb() error {
	var err error
	a.db, err = database.InitDB("postgresql://admin:admin123@localhost:5432/spotnearr")
	if err != nil {
		return err
	}
	if err := database.AutoMigrate(a.db); err != nil {
		return err
	}
	return nil
}
