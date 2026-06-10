package internal_handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	vendorpostgres "github.com/Developer-Aadesh/spotnearr-database/vendordb/postgres"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type InternalHandler struct {
	invProductRepo *vendorpostgres.InvProductRepository
	accessRepo     *vendorpostgres.AccessRepository
}

func NewInternalHandler(db *gorm.DB) *InternalHandler {
	return &InternalHandler{
		invProductRepo: vendorpostgres.NewInvProductRepository(db),
		accessRepo:     vendorpostgres.NewAccessRepository(db),
	}
}

// UserAccessResponse is returned by GET /internal/users/{id}/access.
type UserAccessResponse struct {
	BusinessID uint   `json:"business_id"`
	Role       string `json:"role"`
}

// GetUserAccess handles GET /internal/users/{id}/access.
// Called by User Service at login to embed business_id + role in the JWT.
func (h *InternalHandler) GetUserAccess(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	access, err := h.accessRepo.GetAccessByUserId(context.Background(), id)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UserAccessResponse{
		BusinessID: access.BusinessID,
		Role:       string(access.Role),
	})
}

// GetInventoryProduct handles GET /internal/inventory-products/{id}.
// Called by User Service during claim creation to validate availability and
// freeze a product snapshot at claim time.
func (h *InternalHandler) GetInventoryProduct(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	detail, err := h.invProductRepo.GetByID(context.Background(), id)
	if err != nil {
		if fmt.Sprintf("%v", err) == "record not found" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}

func parseID(r *http.Request, key string) (uint, error) {
	v := strings.TrimSpace(mux.Vars(r)[key])
	if v == "" {
		return 0, fmt.Errorf("missing %s", key)
	}
	id, err := strconv.Atoi(v)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid %s", key)
	}
	return uint(id), nil
}
