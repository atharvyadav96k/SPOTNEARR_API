package get_business_offers

import (
	"errors"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business/get_business_offers/applayer"
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
	if body.BusinessID == (uuid.UUID{}) {
		utils.BadRequest(w, errors.New("business_id is required"))
		return
	}
	app := applayer.Init()
	offers, err := app.GetOffersByBusinessID(body.BusinessID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "offers fetched successfully", models.MapSlice(offers))
}
