package get_reviews

import (
	"errors"
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/users/get_reviews/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
	"github.com/google/uuid"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.BusinessReview](r)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	if body.BusinessID == (uuid.UUID{}) {
		utils.BadRequest(w, errors.New("business_id is required"))
		return
	}
	app := applayer.Init()
	reviews, err := app.GetReviewsByBusinessID(body.BusinessID)
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "reviews fetched successfully", models.MapSlice(reviews))
}
