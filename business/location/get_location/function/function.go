package location

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/location/get_location/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	businessLocations, err := utils.ParseBody[models.BusinessLocation](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	businessLocations, err = app.GetBusinessLocationByID(businessLocations.ID)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, "business location created successfully", businessLocations)
}
