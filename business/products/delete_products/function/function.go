package products_operation

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/products/delete_products/applayer"
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
	err = app.DeleteProduct(product.ID)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.OK(w, "product delete successfully", nil)
}
