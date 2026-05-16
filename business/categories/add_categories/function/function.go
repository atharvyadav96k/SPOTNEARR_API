package categories

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business_add_categories/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	businessCategory, err := utils.ParseBody[models.BusinessCategory](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	businessCategory, err = app.CreateBusinessCategory(*businessCategory)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, "business categories created successfully", businessCategory)
}
