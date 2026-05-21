package applayer

import "github.com/atharvyadav96k/spotnearr-gcp/app"

func Init() *app.App {
	a := &app.App{}
	if err := a.InitDatabase(); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}
	return a
}
