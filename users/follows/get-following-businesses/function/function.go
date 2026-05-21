package follow

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/follow/applayer"
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
	follows, err := app.GetFollowingBusinesses(body.UserID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "following businesses fetched successfully", models.MapSlice(follows))
}
