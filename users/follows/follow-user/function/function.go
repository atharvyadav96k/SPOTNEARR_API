package followuser

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/follow-user/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.UserFollowUser](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	if err := app.FollowUser(body.FollowerID, body.FollowingID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "user followed successfully", nil)
}
