package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
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
	res(w, http.StatusBadGateway, nil)
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
