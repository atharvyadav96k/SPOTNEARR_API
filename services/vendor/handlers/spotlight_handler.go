package handlers

import (
	"net/http"

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
