package users_function

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/create-users/function/applayer"
)

func CreateUsers(w http.ResponseWriter, r *http.Request) {
	app := applayer.Init()

	_ = app
	w.Write([]byte("Users"))
}
