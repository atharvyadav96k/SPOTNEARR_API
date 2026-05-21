package get_liked_businesses

import (
	"errors"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/get_liked_businesses/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
	"github.com/google/uuid"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.UserLikedBusiness](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	if body.UserID == (uuid.UUID{}) {
		utils.BadRequest(w, errors.New("user_id is required"))
		return
	}
	app := applayer.Init()
	likes, err := app.GetLikedBusinesses(body.UserID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "liked businesses fetched successfully", models.MapSlice(likes))
}
