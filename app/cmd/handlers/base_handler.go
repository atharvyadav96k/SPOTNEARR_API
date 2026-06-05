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
	"github.com/atharvyadav96k/SPOTNEARR_API/utils"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/request"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/response"
	"github.com/gorilla/mux"
)

type BaseHandler struct {
	services *services.Services
}

func NewBaserHandler(services *services.Services) *BaseHandler {
	return &BaseHandler{
		services: services,
	}
}

func (b *BaseHandler) GetBizService() *services.BusinessService {
	return b.services.BusinessService
}

func (b *BaseHandler) GetUserService() *services.UserService {
	return b.services.UserService
}

func (b *BaseHandler) GetInvService() *services.InventoryService {
	return b.services.InventoryService
}

func (b *BaseHandler) GetProductService() *services.ProductService {
	return b.services.ProductService
}

func (b *BaseHandler) GetClaimService() *services.ClaimService {
	return b.services.ClaimService
}

func (b *BaseHandler) GetSearchService() *services.SearchService {
	return b.services.SearchService
}

func (b *BaseHandler) GetCategoryService() *services.CategoryService {
	return b.services.CategoryService
}

func (b *BaseHandler) Slug(r *http.Request) string {
	val := request.GetVal(r, "slug")
	if val == nil {
		return ""
	}
	return strings.TrimSpace(val.ToString())
}

func (b *BaseHandler) CategoryIDs(r *http.Request) []uint {
	val := request.GetVal(r, "categoryIds")
	if val == nil {
		return nil
	}
	return val.ToUintSlice()
}

func (b *BaseHandler) StoreIDs(r *http.Request) []uint {
	val := request.GetVal(r, "storeIds")
	if val == nil {
		return nil
	}
	return val.ToUintSlice()
}

func (b *BaseHandler) Email(r *http.Request) (string, error) {
	val := request.GetVal(r, "email")
	if val == nil {
		return "", fmt.Errorf("email field is missing")
	}
	emailStr := val.ToString()
	ok := utils.ValidateEmail(emailStr)
	if !ok {
		return emailStr, fmt.Errorf("Invalid Email")
	}
	return emailStr, nil
}

func (b *BaseHandler) Password(r *http.Request) (string, error) {
	val := request.GetVal(r, "password")
	if val == nil {
		return "", fmt.Errorf("password field is missing")
	}
	passStr := val.ToString()
	ok, message := utils.ValidatePassword(passStr)
	if !ok {
		return passStr, fmt.Errorf("%s", message)
	}
	return passStr, nil
}

func (b *BaseHandler) Phone(r *http.Request) (string, error) {
	val := request.GetVal(r, "phone")
	if val == nil {
		return "", fmt.Errorf("phone field is missing")
	}
	phoneStr := val.ToString()
	ok := utils.ValidatePhone(phoneStr)
	if !ok {
		return phoneStr, fmt.Errorf("Invalid phone number")
	}
	return phoneStr, nil
}

func (b *BaseHandler) Name(r *http.Request) string {
	val := request.GetVal(r, "name")
	if val == nil {
		return ""
	}
	return strings.TrimSpace(val.ToString())
}

func (b *BaseHandler) Address(r *http.Request) string {
	val := request.GetVal(r, "address")
	if val == nil {
		return ""
	}
	return val.ToString()
}

func (b *BaseHandler) Desc(r *http.Request) string {
	val := request.GetVal(r, "desc")
	if val == nil {
		return ""
	}
	return val.ToString()
}

func (b *BaseHandler) Quantity(r *http.Request) *request.ValueUnit {
	val := request.GetVal(r, "quantity")
	if val == nil {
		return nil
	}
	return val.ToValueUnit()
}

func (b *BaseHandler) GeoHash(r *http.Request) string {
	val := request.GetVal(r, "geohash")
	if val == nil {
		return ""
	}
	return val.ToString()
}

func (b *BaseHandler) Lat(r *http.Request) *request.Body {
	val := request.GetVal(r, "lat")
	if val == nil {
		return nil
	}
	return val
}

func (b *BaseHandler) Long(r *http.Request) *request.Body {
	val := request.GetVal(r, "long")
	if val == nil {
		return nil
	}
	return val
}

func (b *BaseHandler) Price(r *http.Request) *request.ValueUnit {
	val := request.GetVal(r, "price")
	if val == nil {
		return nil
	}
	return val.ToValueUnit()
}

func (b *BaseHandler) ProductID(r *http.Request) uint {
	val := request.GetVal(r, "productID")
	if val == nil {
		return 0
	}
	return uint(val.ToInt())
}

func (b *BaseHandler) Count(r *http.Request) *int {
	val := request.GetVal(r, "count")
	if val == nil {
		return nil
	}
	result := val.ToInt()
	return &result
}

func (b *BaseHandler) Available(r *http.Request) *bool {
	val := request.GetVal(r, "available")
	if val == nil {
		return nil
	}
	result := strings.ToLower(val.ToString()) == "true"
	return &result
}

func ParseBody[T any](r *http.Request) (*T, error) {
	var data T

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
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
	res(w, http.StatusBadRequest, response.Res{
		Message: message,
	})
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
	if err != nil {
		return 0
	}
	if claim.BusinessId == nil {
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

func extractKeyFromPath(r *http.Request, key string) (uint, error) {
	vars := mux.Vars(r)
	value := vars[key]

	if strings.TrimSpace(value) == "" {
		return 0, fmt.Errorf("missing or invalid path parameter: %s", key)
	}
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("Invalid Id")
	}
	return uint(id), nil
}

func (h *BaseHandler) GetBusinessId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "bizId")
}

func (h *BaseHandler) GetUserId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "userId")
}

func (h *BaseHandler) GetInventoryId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "invId")
}

func (h *BaseHandler) GetProductId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "productId")
}

func (h *BaseHandler) GetOfferId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "offerId")
}

func (h *BaseHandler) GetSpotlightId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "spotlightId")
}

func (h *BaseHandler) GetInvProductId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "invProductId")
}

func (h *BaseHandler) GetClaimId(r *http.Request) (uint, error) {
	return extractKeyFromPath(r, "claimId")
}

func (h *BaseHandler) QuerySession(r *http.Request) string {
	return r.URL.Query().Get("session")
}
