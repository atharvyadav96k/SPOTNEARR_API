package comments

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/spotlight/comments/create/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	comments, err := utils.ParseAndValidate[models.SpotlightComment](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	spotlight, err := app.CreateSpotlightComment(*comments)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "spotlight comment created successfully", spotlight.ToResponse())
}
