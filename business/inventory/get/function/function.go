package inventory

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/inventory/get/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	inventory, err := utils.ParseBody[models.ProductInventory](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	inventory, err = app.GetProductInventoryByID(inventory.ID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "inventry", inventory)
}
