package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/dtos"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/models"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
)

type InventoryHandler struct {
	BaseHandler
	svc *services.InventoryService
}

func NewInventoryHandler(svc *services.InventoryService) *InventoryHandler {
	return &InventoryHandler{svc: svc}
}

func (i *InventoryHandler) InventoryCreate(w http.ResponseWriter, r *http.Request) {
	bizID := i.ClaimGetBusinessID(r)
	if bizID == 0 {
		i.ResponseBadRequest(w)
		return
	}
	var dto dtos.InventoryCreateRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		i.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	store := models.NewStore(dto.Name, dto.StreetAddress, dto.Lat, dto.Long)
	i.Response(w, i.svc.CreateInventory(bizID, &store))
}

func (i *InventoryHandler) InventoryUpdate(w http.ResponseWriter, r *http.Request) {
	bizID := i.ClaimGetBusinessID(r)
	if bizID == 0 {
		i.ResponseBadRequest(w)
		return
	}
	storeID, err := i.GetInventoryID(r)
	if err != nil {
		i.ResponseBadRequest(w)
		return
	}
	var dto dtos.InventoryUpdateRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		i.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	i.Response(w, i.svc.UpdateInventory(bizID, storeID, dto.Name, dto.Address, *dto.Lat, *dto.Long))
}

func (i *InventoryHandler) InventoryGetProducts(w http.ResponseWriter, r *http.Request) {
	bizID := i.ClaimGetBusinessID(r)
	invID, err := i.GetInventoryID(r)
	if err != nil {
		i.ResponseBadRequest(w)
		return
	}
	i.Response(w, i.svc.GetInvProducts(bizID, invID))
}

func (i *InventoryHandler) InventoryAddProduct(w http.ResponseWriter, r *http.Request) {
	bizID := i.ClaimGetBusinessID(r)
	invID, err := i.GetInventoryID(r)
	if err != nil {
		i.ResponseBadRequest(w)
		return
	}
	var dto dtos.InventoryAddProductRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		i.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	i.Response(w, i.svc.AddInvProduct(bizID, invID, dto.ProductID, dto.Count, dto.Available))
}

func (i *InventoryHandler) InventoryUpdateProduct(w http.ResponseWriter, r *http.Request) {
	bizID := i.ClaimGetBusinessID(r)
	invProdID, err := i.GetInvProductID(r)
	if err != nil {
		i.ResponseBadRequest(w)
		return
	}
	var dto dtos.InventoryUpdateProductRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		i.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	i.Response(w, i.svc.UpdateInvProduct(bizID, invProdID, dto.Count, dto.Available))
}

func (i *InventoryHandler) InventoryRemoveProduct(w http.ResponseWriter, r *http.Request) {
	bizID := i.ClaimGetBusinessID(r)
	invProdID, err := i.GetInvProductID(r)
	if err != nil {
		i.ResponseBadRequest(w)
		return
	}
	i.Response(w, i.svc.RemoveInvProduct(bizID, invProdID))
}
