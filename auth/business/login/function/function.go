package auth

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business_login/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.User](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	user, err := app.GetUserByEmail(*body.Email)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	if user.Role != "business" {
		utils.Unauthorized(w, "not a business account")
		return
	}
	if !user.IsVerified {
		utils.Unauthorized(w, "account is not verified")
		return
	}
	utils.OK(w, "business login successful", user.ToResponse())
}
