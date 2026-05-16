package categories

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/business_get_categories/applayer"
	"github.com/atharvyadav96k/spotnearr-gcp/app/utils"
)

func Function(w http.ResponseWriter, r *http.Request) {
	app := applayer.Init()
	categories, err := app.GetBusinessCategories()
	if err != nil {
		utils.BadRequest(w, err)
		return
	}
	utils.OK(w, "business categories loaded successfully", categories)
}
