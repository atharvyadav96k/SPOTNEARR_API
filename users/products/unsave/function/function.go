package unsaveproduct

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/unsave-product/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseAndValidate[models.UserSavedProduct](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	if err := app.UnsaveProduct(body.UserID, body.ProductID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "product unsaved successfully", nil)
}
