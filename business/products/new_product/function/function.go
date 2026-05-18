package products_operation

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/products/new_product/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	product, err := utils.ParseBody[models.Product](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	product, err = app.CreateProduct(*product)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, "product created successfully", product)
}
