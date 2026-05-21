package delete_offer

import (
	"errors"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business/delete_offer/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
	"github.com/google/uuid"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.Offer](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	if body.ID == (uuid.UUID{}) {
		utils.BadRequest(w, errors.New("id is required"))
		return
	}
	if body.BusinessID == (uuid.UUID{}) {
		utils.BadRequest(w, errors.New("business_id is required"))
		return
	}
	app := applayer.Init()
	if err := app.DeleteOffer(body.ID, body.BusinessID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "offer deleted successfully", nil)
}
