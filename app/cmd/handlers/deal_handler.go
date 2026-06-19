package handlers

import (
	"net/http"
	"strconv"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type DealHandler struct{ BaseHandler }

func NewDealHandler(svcs *services.Services) *DealHandler {
	return &DealHandler{BaseHandler: *NewBaserHandler(svcs)}
}

func (h *DealHandler) NearbyDeals(w http.ResponseWriter, r *http.Request) {
	lat, err1 := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
	long, err2 := strconv.ParseFloat(r.URL.Query().Get("long"), 64)
	if err1 != nil || err2 != nil || lat < -90 || lat > 90 || long < -180 || long > 180 {
		h.ResponseBadRequestWithMessage(w, "invalid lat/long parameters")
		return
	}
	h.Response(w, h.services.DealService.NearbyDeals(lat, long))
}
