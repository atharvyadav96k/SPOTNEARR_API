package auth

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business_register/applayer"
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
	business, err = app.RegisterBusiness(*business)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, "business registered successfully", business)

}
