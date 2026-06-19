package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
	userdb "github.com/Developer-Aadesh/spotnearr-database/user"
	"github.com/gorilla/mux"
)

type InternalHandler struct {
	userService  *services.UserService
	claimService *services.ClaimService
}

func NewInternalHandler(s *services.Services) *InternalHandler {
	return &InternalHandler{
		userService:  s.UserService,
		claimService: s.ClaimService,
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (h *InternalHandler) InvalidateRefresh(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(mux.Vars(r)["userId"])
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid userId: %s", idStr), http.StatusBadRequest)
		return
	}
	if err := h.userService.InvalidateRefreshToken(uint(id)); err != nil {
		log.Printf("internal: invalidate-refresh user %d: %v", id, err)
	}
	w.WriteHeader(http.StatusOK)
}

// GetClaimsByProductIDs handles GET /internal/claims?product_ids=1,2,3
func (h *InternalHandler) GetClaimsByProductIDs(w http.ResponseWriter, r *http.Request) {
	raw := strings.TrimSpace(r.URL.Query().Get("product_ids"))
	if raw == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "product_ids required"})
		return
	}
	parts := strings.Split(raw, ",")
	var ids []uint
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || n <= 0 {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid product_id: " + p})
			return
		}
		ids = append(ids, uint(n))
	}
	claims, err := h.claimService.GetClaimsByProductIDs(r.Context(), ids)
	if err != nil {
		log.Printf("internal: get claims by product ids: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	writeJSON(w, http.StatusOK, claims)
}

// GetClaimByID handles GET /internal/claims/{claimId}
func (h *InternalHandler) GetClaimByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(mux.Vars(r)["claimId"])
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid claimId"})
		return
	}
	claim, err := h.claimService.GetClaimByID(r.Context(), uint(id))
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, claim)
}

// UpdateClaimStatus handles PATCH /internal/claims/{claimId}/status
func (h *InternalHandler) UpdateClaimStatus(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(mux.Vars(r)["claimId"])
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid claimId"})
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	status := userdb.ClaimStatus(body.Status)
	if status != userdb.ClaimStatusAccepted && status != userdb.ClaimStatusRejected {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "status must be 'accepted' or 'rejected'"})
		return
	}
	if err := h.claimService.UpdateClaimStatus(r.Context(), uint(id), status); err != nil {
		log.Printf("internal: update claim %d status: %v", id, err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
		return
	}
	w.WriteHeader(http.StatusOK)
}
