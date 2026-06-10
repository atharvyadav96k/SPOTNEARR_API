package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	pkgmid "github.com/atharvyadav96k/spotnearr/pkg/middleware"
	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/pkg/jwtutil"
	"github.com/gorilla/mux"
)

type BaseHandler struct{}

func respond(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(body)
}

func (h *BaseHandler) Response(w http.ResponseWriter, r httputil.Res) {
	respond(w, r.StatusCode, r)
}

func (h *BaseHandler) ResponseBadRequest(w http.ResponseWriter) {
	respond(w, http.StatusBadRequest, nil)
}

func (h *BaseHandler) ResponseBadRequestWithMessage(w http.ResponseWriter, message string) {
	respond(w, http.StatusBadRequest, httputil.Res{Message: message})
}

func (h *BaseHandler) ResponseOK(w http.ResponseWriter) {
	respond(w, http.StatusOK, nil)
}

// Validatable is implemented by any DTO that can validate itself.
type Validatable interface {
	Validate() error
}

// parseAndValidateBody decodes JSON into dst then calls Validate().
func parseAndValidateBody(r *http.Request, dst Validatable) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid request body")
	}
	return dst.Validate()
}

func ParseBody[T any](r *http.Request) (*T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (h *BaseHandler) getClaims(r *http.Request) (*jwtutil.UserClaims, error) {
	claims, ok := pkgmid.ClaimsFromContext(r)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}

func (h *BaseHandler) ClaimGetUserID(r *http.Request) uint {
	claims, err := h.getClaims(r)
	if err != nil {
		return 0
	}
	return claims.UserId
}

func (h *BaseHandler) ClaimGetBusinessID(r *http.Request) uint {
	claims, err := h.getClaims(r)
	if err != nil || claims.BusinessId == nil {
		return 0
	}
	return *claims.BusinessId
}

func extractPath(r *http.Request, key string) (uint, error) {
	vars := mux.Vars(r)
	v := vars[key]
	if strings.TrimSpace(v) == "" {
		return 0, fmt.Errorf("missing path parameter: %s", key)
	}
	id, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}
	return uint(id), nil
}

func (h *BaseHandler) GetBizID(r *http.Request) (uint, error)        { return extractPath(r, "bizId") }
func (h *BaseHandler) GetInventoryID(r *http.Request) (uint, error)  { return extractPath(r, "invId") }
func (h *BaseHandler) GetProductID(r *http.Request) (uint, error)    { return extractPath(r, "productId") }
func (h *BaseHandler) GetInvProductID(r *http.Request) (uint, error) { return extractPath(r, "invProductId") }
