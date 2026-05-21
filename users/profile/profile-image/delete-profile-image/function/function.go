package profile_image

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/delete_profile_image/applayer"
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
	if err := app.DeleteUserAvatar(body.ID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "profile image deleted successfully", nil)
}
