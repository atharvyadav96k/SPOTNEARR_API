package products

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/nearby_products/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.BusinessLocation](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	products, err := app.GetNearbyProducts(body.Latitude, body.Longitude, 5)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "nearby products fetched successfully", models.MapSlice(products))
}
