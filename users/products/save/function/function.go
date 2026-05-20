package saveproduct

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/save-product/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.UserSavedProduct](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	if err := app.SaveProduct(*body); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "product saved successfully", nil)
}
