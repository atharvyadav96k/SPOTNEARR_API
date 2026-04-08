package applayer

import "github.com/atharvyadav96k/SPOTNEARR_SHARED/app"

func Init() *application {
	app := app.Init()

	return &application{
		App: app,
	}
}
