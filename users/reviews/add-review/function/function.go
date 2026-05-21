package add_review

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/add_review/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	review, err := utils.ParseAndValidate[models.BusinessReview](r)
	if err != nil {
		utils.ValidationError(w, err)
		return
	}
	app := applayer.Init()
	created, err := app.CreateReview(*review)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.Created(w, "review submitted successfully", created.ToResponse())
}
