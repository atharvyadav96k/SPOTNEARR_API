package getlikedproducts

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/get-liked-products/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.UserLikedProduct](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	app := applayer.Init()
	liked, err := app.GetLikedProducts(body.UserID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "liked products fetched successfully", liked)
}
