package getsavedspotlights

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/get-saved-spotlights/applayer"
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
	spotlights, err := app.GetSavedSpotlightsByUserID(body.ID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "saved spotlights retrieved successfully", models.MapSlice(spotlights))
}
