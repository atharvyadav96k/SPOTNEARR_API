package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/atharvyadav96k/spotnearr/pkg/httputil"
	"github.com/atharvyadav96k/spotnearr/pkg/jwtutil"
	middleware "github.com/atharvyadav96k/spotnearr/vendor-svc/middlewares"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/utils/request"
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

func (h *BaseHandler) getClaims(r *http.Request) (*jwtutil.UserClaims, error) {
	claims, ok := r.Context().Value(middleware.ClaimsKey).(*jwtutil.UserClaims)
	if !ok || claims == nil {
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

func (h *BaseHandler) Name(r *http.Request) string {
	v := request.GetVal(r, "name")
	if v == nil {
		return ""
	}
	return strings.TrimSpace(v.ToString())
}

func (h *BaseHandler) Desc(r *http.Request) string {
	v := request.GetVal(r, "desc")
	if v == nil {
		return ""
	}
	return v.ToString()
}

func (h *BaseHandler) Slug(r *http.Request) string {
	v := request.GetVal(r, "slug")
	if v == nil {
		return ""
	}
	return strings.TrimSpace(v.ToString())
}

func (h *BaseHandler) Address(r *http.Request) string {
	v := request.GetVal(r, "address")
	if v == nil {
		return ""
	}
	return v.ToString()
}

func (h *BaseHandler) Lat(r *http.Request) *request.Body {
	return request.GetVal(r, "lat")
}

func (h *BaseHandler) Long(r *http.Request) *request.Body {
	return request.GetVal(r, "long")
}

func (h *BaseHandler) Price(r *http.Request) *request.ValueUnit {
	v := request.GetVal(r, "price")
	if v == nil {
		return nil
	}
	return v.ToValueUnit()
}

func (h *BaseHandler) Quantity(r *http.Request) *request.ValueUnit {
	v := request.GetVal(r, "quantity")
	if v == nil {
		return nil
	}
	return v.ToValueUnit()
}

func (h *BaseHandler) CategoryIDs(r *http.Request) []uint {
	v := request.GetVal(r, "categoryIds")
	if v == nil {
		return nil
	}
	return v.ToUintSlice()
}

func (h *BaseHandler) StoreIDs(r *http.Request) []uint {
	v := request.GetVal(r, "storeIds")
	if v == nil {
		return nil
	}
	return v.ToUintSlice()
}

func (h *BaseHandler) Count(r *http.Request) *int {
	v := request.GetVal(r, "count")
	if v == nil {
		return nil
	}
	n := v.ToInt()
	return &n
}

func (h *BaseHandler) Available(r *http.Request) *bool {
	v := request.GetVal(r, "available")
	if v == nil {
		return nil
	}
	b := strings.ToLower(v.ToString()) == "true"
	return &b
}

func (h *BaseHandler) ProductIDField(r *http.Request) uint {
	v := request.GetVal(r, "productID")
	if v == nil {
		return 0
	}
	return uint(v.ToInt())
}

func ParseBody[T any](r *http.Request) (*T, error) {
	var data T
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}
