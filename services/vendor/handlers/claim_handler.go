package handlers

import (
	"net/http"

	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
)

type ClaimHandler struct {
	BaseHandler
	svc *services.ClaimService
}

func NewClaimHandler(svc *services.ClaimService) *ClaimHandler {
	return &ClaimHandler{svc: svc}
}

func (h *ClaimHandler) ListClaims(w http.ResponseWriter, r *http.Request) {
	bizID := h.ClaimGetBusinessID(r)
	h.Response(w, h.svc.ListClaims(bizID))
}

func (h *ClaimHandler) AcceptClaim(w http.ResponseWriter, r *http.Request) {
	h.updateStatus(w, r, "accepted")
}

func (h *ClaimHandler) RejectClaim(w http.ResponseWriter, r *http.Request) {
	h.updateStatus(w, r, "rejected")
}

func (h *ClaimHandler) updateStatus(w http.ResponseWriter, r *http.Request, status string) {
	bizID := h.ClaimGetBusinessID(r)
	claimID, err := extractPath(r, "claimId")
	if err != nil {
		h.ResponseBadRequestWithMessage(w, "invalid claim ID")
		return
	}
	dto := pkgdtos.ClaimStatusUpdateRequest{Status: status}
	h.Response(w, h.svc.UpdateClaimStatus(bizID, claimID, dto))
}
