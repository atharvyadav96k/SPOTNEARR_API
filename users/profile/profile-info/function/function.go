package users

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/profile/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.User](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	user, err := app.GetUserByID(body.ID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "user profile fetched successfully", user.ToResponse())
}
