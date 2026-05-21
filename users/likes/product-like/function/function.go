package productlike

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/user/product-like/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseAndValidate[models.UserLikedProduct](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	if err := app.LikeProduct(body.UserID, body.ProductID); err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "product liked successfully", nil)
}
