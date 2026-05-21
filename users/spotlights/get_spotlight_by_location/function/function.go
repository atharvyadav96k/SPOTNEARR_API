package spotlights

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/spotlight_get_by_location/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.BusinessLocation](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	spotlights, err := app.GetSpotlightsByLocation(body.Latitude, body.Longitude)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "spotlights fetched successfully", models.MapSlice(spotlights))
}
