package products_operation

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/products/get_business_prouducts/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	business, err := utils.ParseBody[models.Business](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	products, err := app.GetProductsByBusinessID(business.ID)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.OK(w, "products fetched successfully", models.MapSlice(products))
}
