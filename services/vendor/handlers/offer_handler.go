package handlers

import (
	"net/http"

	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
)

type OfferHandler struct {
	BaseHandler
	svc *services.OfferService
}

func NewOfferHandler(svc *services.OfferService) *OfferHandler {
	return &OfferHandler{svc: svc}
}

func (h *OfferHandler) OfferCreate(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	var dto pkgdtos.OfferCreateRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.svc.CreateOffer(bizID, dto))
}

func (h *OfferHandler) OfferList(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	h.Response(w, h.svc.ListOffers(bizID))
}

func (h *OfferHandler) OfferGet(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	offerID, err := extractPath(r, "offerId")
	if err != nil {
		h.ResponseBadRequestWithMessage(w, "invalid offer ID")
		return
	}
	h.Response(w, h.svc.GetOffer(bizID, offerID))
}

func (h *OfferHandler) OfferUpdate(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	offerID, err := extractPath(r, "offerId")
	if err != nil {
		h.ResponseBadRequestWithMessage(w, "invalid offer ID")
		return
	}
	var dto pkgdtos.OfferUpdateRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.svc.UpdateOffer(bizID, offerID, dto))
}

func (h *OfferHandler) OfferDelete(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	offerID, err := extractPath(r, "offerId")
	if err != nil {
		h.ResponseBadRequestWithMessage(w, "invalid offer ID")
		return
	}
	h.Response(w, h.svc.DeleteOffer(bizID, offerID))
}
