package handlers

import (
	"net/http"
	"strconv"

	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
)

type SpotlightHandler struct {
	BaseHandler
	svc *services.SpotlightService
}

func NewSpotlightHandler(svc *services.SpotlightService) *SpotlightHandler {
	return &SpotlightHandler{svc: svc}
}

func (h *SpotlightHandler) SpotlightPost(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	var dto pkgdtos.SpotlightCreateRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.svc.PostSpotlight(bizID, dto))
}

func (h *SpotlightHandler) SpotlightList(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	h.Response(w, h.svc.ListSpotlights(bizID))
}

func (h *SpotlightHandler) SpotlightGet(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	id, err := extractPath(r, "spotlightId")
	if err != nil {
		h.ResponseBadRequestWithMessage(w, "invalid spotlight ID")
		return
	}
	h.Response(w, h.svc.GetSpotlight(bizID, id))
}

func (h *SpotlightHandler) SpotlightDelete(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	id, err := extractPath(r, "spotlightId")
	if err != nil {
		h.ResponseBadRequestWithMessage(w, "invalid spotlight ID")
		return
	}
	h.Response(w, h.svc.DeleteSpotlight(bizID, id))
}

func (h *SpotlightHandler) SpotlightFeed(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	lat, errLat := strconv.ParseFloat(q.Get("lat"), 64)
	long, errLong := strconv.ParseFloat(q.Get("long"), 64)
	if errLat != nil || errLong != nil {
		h.ResponseBadRequestWithMessage(w, "lat and long are required")
		return
	}
	rangeKm := 10.0
	if v, err := strconv.ParseFloat(q.Get("range"), 64); err == nil && v > 0 {
		rangeKm = v
	}
	h.Response(w, h.svc.GetFeed(lat, long, rangeKm))
}
