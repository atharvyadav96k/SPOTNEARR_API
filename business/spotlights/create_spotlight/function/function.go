package spotlight_operation

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business/spotlight_create/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	spotlight, err := utils.ParseAndValidate[models.Spotlight](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	spotlight, err = app.CreateSpotlight(*spotlight)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "spotlight created successfully", spotlight.ToResponse())
}
