package create_offer

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business/create_offer/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	offer, err := utils.ParseAndValidate[models.Offer](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	created, err := app.CreateOffer(*offer)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "offer created successfully", created.ToResponse())
}
