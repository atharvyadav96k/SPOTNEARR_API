package getfollowingusers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/get-following-users/applayer"
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
	following, err := app.GetUserFollowing(body.FollowerID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "following retrieved successfully", models.MapSlice(following))
}
