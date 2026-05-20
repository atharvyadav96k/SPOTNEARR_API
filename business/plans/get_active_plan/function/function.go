package plans

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business_active_plan/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[models.Business](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	sub, err := app.GetBusinessActivePlan(body.ID)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.OK(w, "active plan fetched successfully", sub)
}
