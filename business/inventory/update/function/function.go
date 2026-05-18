package inventory

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/inventory/update/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	productInventory, err := utils.ParseBody[models.ProductInventory](r)
	if err != nil {
		utils.BadRequest(w, err)
	}
	app := applayer.Init()
	productInventory, err = app.UpdateProductInventory(*productInventory)
	if err != nil {
		utils.BadRequest(w, err)
	}
	utils.Created(w, "product updated successfully", productInventory)
}
