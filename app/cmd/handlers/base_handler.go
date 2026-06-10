package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/auth"
	middleware "github.com/atharvyadav96k/SPOTNEARR_API/middlewares"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"github.com/gorilla/mux"
)

type BaseHandler struct {
	services *services.Services
}

func NewBaserHandler(services *services.Services) *BaseHandler {
	return &BaseHandler{services: services}
}

func (b *BaseHandler) GetUserService() *services.UserService {
	return b.services.UserService
}

func (b *BaseHandler) GetClaimService() *services.ClaimService {
	return b.services.ClaimService
}

func (b *BaseHandler) GetReviewService() *services.ReviewService {
	return b.services.ReviewService
}

// ParseBody decodes the JSON request body into T.
func ParseBody[T any](r *http.Request) (*T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

// Validatable is implemented by any DTO that can validate itself.
type Validatable interface {
	Validate() error
}

// parseAndValidateBody decodes JSON into dst (must be a pointer to a DTO) then
// calls Validate(). Returns a single error covering both steps.
func parseAndValidateBody(r *http.Request, dst Validatable) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid request body")
	}
	return dst.Validate()
}

func res(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func (h *BaseHandler) ResponseOK(w http.ResponseWriter) {
	res(w, http.StatusOK, nil)
}

func (h *BaseHandler) Response(w http.ResponseWriter, r response.Res) {
	res(w, r.StatusCode, r)
}

func (h *BaseHandler) ResponseBadRequest(w http.ResponseWriter) {
	res(w, http.StatusBadRequest, nil)
}

func (h *BaseHandler) ResponseNotFound(w http.ResponseWriter) {
	res(w, http.StatusNotFound, nil)
}

func (h *BaseHandler) ResponseBadRequestWithMessage(w http.ResponseWriter, message string) {
	res(w, http.StatusBadRequest, response.Res{Message: message})
}

func (h *BaseHandler) getClaims(r *http.Request) (auth.UserClaims, error) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*auth.UserClaims)
	if !ok {
		return auth.UserClaims{}, fmt.Errorf("unauthorized access")
	}
	return *claims, nil
}

func (h *BaseHandler) ClaimGetBusinessId(r *http.Request) uint {
	claim, err := h.getClaims(r)
	if err != nil || claim.BusinessId == nil {
		return 0
	}
	return *claim.BusinessId
}

func (h *BaseHandler) ClaimGetUserId(r *http.Request) uint {
	claim, err := h.getClaims(r)
	if err != nil {
		return 0
	}
	return claim.UserId
}

func (h *BaseHandler) QuerySession(r *http.Request) string {
	return r.URL.Query().Get("session")
}

func extractKeyFromPath(r *http.Request, key string) (uint, error) {
	vars := mux.Vars(r)
	value := vars[key]
	if strings.TrimSpace(value) == "" {
		return 0, fmt.Errorf("missing or invalid path parameter: %s", key)
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return uint(id), nil
}

func (h *BaseHandler) GetUserId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "userId")
}

func (h *BaseHandler) GetInvProductId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "invProductId")
}

func (h *BaseHandler) GetClaimId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "claimId")
}

func (h *BaseHandler) GetOfferId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "offerId")
}

func (h *BaseHandler) GetSpotlightId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "spotlightId")
}
