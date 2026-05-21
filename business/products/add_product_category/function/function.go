package add_product_category

import (
	"encoding/json"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business_add_product_category/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

type request struct {
	Name string `json:"name"`
}

func Function(w http.ResponseWriter, r *http.Request) {
	var body request
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		utils.ValidationError(w, "name is required")
		return
	}
	app := applayer.Init()
	cat, err := app.CreateProductCategory(body.Name)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, "product category created successfully", cat.ToResponse())
}
