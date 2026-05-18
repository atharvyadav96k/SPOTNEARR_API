package inventory

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/inventory/get_by_business/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	business, err := utils.ParseBody[models.Business](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}

	app := applayer.Init()
	inventories, err := app.GetProductInventoriesByBusinessID(business.ID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "product inventories by business", inventories)
}
