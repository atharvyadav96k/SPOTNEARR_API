package applayer

import "github.com/atharvyadav96k/spotnearr-gcp/app"

func Init() *app.App {
	a := app.App{}
	if err := a.InitDatabase(); err != nil {
		panic(err)
	}
	return &a
}
