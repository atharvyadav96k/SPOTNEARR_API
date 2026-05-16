package auth

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user_login/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	user, err := utils.ParseBody[models.User](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	user, err = app.GetUserByEmail(*user.Email)
	if err != nil {
		utils.Conflict(w, err.Error())
		return
	}
	if !user.IsVerified {
		utils.Unauthorized(w, "user is not verified")
		return
	}
	utils.OK(w, "user login successful", user)
}
