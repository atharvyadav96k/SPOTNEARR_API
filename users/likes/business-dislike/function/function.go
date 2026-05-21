package user_likes

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/like/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseAndValidate[models.UserLikedBusiness](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	if err := app.UnlikeBusiness(body.UserID, body.BusinessID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "business disliked successfully", nil)
}
