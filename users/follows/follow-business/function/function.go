package follow

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/follow/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseAndValidate[models.UserFollowBusiness](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	if err := app.FollowBusiness(body.UserID, body.BusinessID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "business followed successfully", nil)
}
