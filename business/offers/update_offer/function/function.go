package update_offer

import (
	"errors"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business/update_offer/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
	"github.com/google/uuid"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseAndValidate[models.Offer](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	if body.ID == (uuid.UUID{}) {
		utils.BadRequest(w, errors.New("id is required"))
		return
	}
	app := applayer.Init()
	updated, err := app.UpdateOffer(*body)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "offer updated successfully", updated.ToResponse())
}
