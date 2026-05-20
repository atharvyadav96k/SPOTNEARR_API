package plans

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business_set_active_plan/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/database/models"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

// SetPlanRequest binds the incoming JSON body for activating a plan.
type SetPlanRequest struct {
	BusinessID models.Business `json:"business"`
	PlanID     models.Plan     `json:"plan"`
}

func Function(w http.ResponseWriter, r *http.Request) {
	body, err := utils.ParseBody[SetPlanRequest](r)
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	app := applayer.Init()
	if err := app.SetBusinessActivePlan(body.BusinessID.ID, body.PlanID.ID); err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.OK(w, "plan activated successfully", nil)
}
