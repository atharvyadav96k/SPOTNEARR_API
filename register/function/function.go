package demofunction

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/register/applayer"
	"github.com/atharvyadav96k/SPOTNEARR_SHARED/common/res"
)

func Register(w http.ResponseWriter, r *http.Request) {
	app := applayer.Init()
	if app == nil {
		panic("Failed to initialize the application")
	}
	res.Send(w, "201", "User registered successfully", nil)
}
