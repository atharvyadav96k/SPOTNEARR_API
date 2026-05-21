package get_notifications

import (
	"errors"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/get_notifications/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
	"github.com/google/uuid"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.Notification](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	if body.UserID == (uuid.UUID{}) {
		utils.BadRequest(w, errors.New("user_id is required"))
		return
	}
	app := applayer.Init()
	notifications, err := app.GetNotificationsByUserID(body.UserID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "notifications fetched successfully", models.MapSlice(notifications))
}
