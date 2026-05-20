package plans

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business_plans_add/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	plan, err := utils.ParseBody[models.Plan](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	created, err := app.AddPlan(*plan)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.Created(w, "plan added successfully", created)
}
