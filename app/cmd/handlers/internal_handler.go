package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
	"github.com/gorilla/mux"
)

type InternalHandler struct {
	userService *services.UserService
}

func NewInternalHandler(s *services.Services) *InternalHandler {
	return &InternalHandler{userService: s.UserService}
}

// InvalidateRefresh handles POST /internal/users/{id}/invalidate-refresh.
// Called by the Vendor Service (fire-and-forget) after business registration so
// the user's next login generates a JWT that contains business_id.
func (h *InternalHandler) InvalidateRefresh(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(mux.Vars(r)["userId"])
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, fmt.Sprintf("invalid userId: %s", idStr), http.StatusBadRequest)
		return
	}
	if err := h.userService.InvalidateRefreshToken(uint(id)); err != nil {
		log.Printf("internal: invalidate-refresh user %d: %v", id, err)
		// Still return 200 — this is fire-and-forget; the vendor side ignores failures.
	}
	w.WriteHeader(http.StatusOK)
}
