package comments

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/spotlight/comments/get_by_spotlight/applayer"
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
	comments, err := app.GetSpotlightCommentsBySpotlightID(spotlight.ID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "spotlight comment created successfully", comments)
}
