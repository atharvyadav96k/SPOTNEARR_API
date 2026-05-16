package auth

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
	"github.com/atharvyadav96k/spotnearr-gcp/auth/user/register/applayer"
)

func Function(
	w http.ResponseWriter,
	r *http.Request,
) {

	app := applayer.Init()

	user, err := utils.ParseBody[models.User](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}

	err = app.RegisterUser(*user)
	if err != nil {
		utils.Conflict(w, err.Error())
		return
	}

	utils.Created(w, "user created successfully", user)
}
