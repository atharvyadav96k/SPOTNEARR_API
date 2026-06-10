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
