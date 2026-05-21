package inventory

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/inventory/get_by_business/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
	"github.com/google/uuid"
)

type getByBusinessRequest struct {
	BusinessID uuid.UUID `json:"business_id"`
}

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[getByBusinessRequest](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	app := applayer.Init()
	inventories, err := app.GetProductInventoriesByBusinessID(body.BusinessID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "product inventories by business", inventories)
}
