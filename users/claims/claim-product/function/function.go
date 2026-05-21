package claim_product

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/claim_product/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseAndValidate[models.ProductClaim](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	claim, err := app.CreateClaim(body.UserID, body.BusinessID, body.InventoryID, body.Quantity, body.Note)
	if err != nil {
		if err.Error() == "out of stock" {
			utils.Conflict(w, err.Error())
			return
		}
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "product claimed successfully", claim.ToResponse())
}
