package profile_image

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business/image/upload/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func SetProfileImageInDb(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.Business](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	if body.LogoURL == nil || *body.LogoURL == "" {
		utils.ValidationError(w, "logo_url is required")
		return
	}
	app := applayer.Init()
	if err := app.UpdateBusinessLogo(body.ID, *body.LogoURL); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "profile image set successfully", nil)
}
