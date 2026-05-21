package inventory

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/inventory/create/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	inventorys, err := utils.ParseBody[[]models.ProductInventory](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	if err := utils.ValidateSlice(*inventorys); err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	err = app.CreateProductInventories(*inventorys)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "product inventory created successfully", nil)
}
