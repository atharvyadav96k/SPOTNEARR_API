package products

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/search_products/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

type searchRequest struct {
	Query string `json:"query"`
}

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[searchRequest](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	if body.Query == "" {
		utils.ValidationError(w, "query is required")
		return
	}
	app := applayer.Init()
	products, err := app.SearchProducts(body.Query)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "products fetched successfully", models.MapSlice(products))
}
