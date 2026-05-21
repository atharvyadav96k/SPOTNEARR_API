package profile_image

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business/image/download/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func DeleteProfileImage(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.Business](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	if err := app.DeleteBusinessLogo(body.ID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "profile image deleted successfully", nil)
}
