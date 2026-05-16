package applayer

import (
	"github.com/atharvyadav96k/spotnearr-gcp/app"
)

func Init() *app.App {
	app := &app.App{}
	if err := app.InitDatabase(); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}
	return app
}
