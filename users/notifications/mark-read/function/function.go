package mark_read

import (
	"errors"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/mark_notification_read/applayer"
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
	if body.ID != (uuid.UUID{}) {
		if err := app.MarkNotificationRead(body.ID, body.UserID); err != nil {
			utils.BadRequest(w, err)
			return
		}
		utils.OK(w, "notification marked as read", nil)
		return
	}
	if err := app.MarkAllNotificationsRead(body.UserID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "all notifications marked as read", nil)
}
