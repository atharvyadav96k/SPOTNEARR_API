package plans

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business_plans/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	app := applayer.Init()
	plans, err := app.GetPlans()
	if err != nil {
		utils.BadRequest(w, err.Error())
		return
	}
	utils.OK(w, "plans fetched successfully", plans)
}
