package unfollow

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/unfollow/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.UserFollowBusiness](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	if err := app.UnfollowBusiness(body.UserID, body.BusinessID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "business unfollowed successfully", nil)
}
