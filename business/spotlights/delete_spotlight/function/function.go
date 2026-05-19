package spotlight_operation

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/spotlight/delete/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	spotlight, err := utils.ParseBody[models.Spotlight](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	err = app.DeleteSpotlight(spotlight.ID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "deleted successfully", nil)
}
