package location

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/location/create_location/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	businessLocation, err := utils.ParseBody[models.BusinessLocation](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	businessLocation, err = app.CreateBusinessLocation(*businessLocation)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, "business location created successfully", businessLocation)
}
