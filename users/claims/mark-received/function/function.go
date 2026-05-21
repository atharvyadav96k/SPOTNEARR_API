package mark_received

import (
	"errors"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/mark_received/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
	"github.com/google/uuid"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.ProductClaim](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	if body.ID == (uuid.UUID{}) {
		utils.BadRequest(w, errors.New("id is required"))
		return
	}
	if body.UserID == (uuid.UUID{}) {
		utils.BadRequest(w, errors.New("user_id is required"))
		return
	}
	app := applayer.Init()
	claim, err := app.GetClaimByID(body.ID)
	if err != nil {
		utils.NotFound(w, errors.New("claim not found"))
		return
	}
	if claim.UserID != body.UserID {
		utils.Forbidden(w, errors.New("not your claim"))
		return
	}
	if claim.Status != models.ClaimAccepted {
		utils.BadRequest(w, errors.New("only accepted claims can be marked as received"))
		return
	}
	if err := app.SetClaimStatus(body.ID, models.ClaimCompleted); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "claim marked as received", nil)
}
