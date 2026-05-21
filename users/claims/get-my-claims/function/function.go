package get_my_claims

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/get_my_claims/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.ProductClaim](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	claims, err := app.GetClaimsByUserID(body.UserID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "claims fetched successfully", models.MapSlice(claims))
}
