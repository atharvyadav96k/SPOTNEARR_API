package handlers

import (
	"net/http"

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
	store, err := ParseBody[models.Store](r)
	if err != nil || store == nil {
		i.ResponseBadRequest(w)
		return
	}
	i.Response(w, i.svc.CreateInventory(bizID, store))
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
	latBody := i.Lat(r)
	longBody := i.Long(r)
	if latBody == nil || longBody == nil {
		i.ResponseBadRequestWithMessage(w, "lat and long are required")
		return
	}
	i.Response(w, i.svc.UpdateInventory(bizID, storeID, i.Name(r), i.Address(r), latBody.ToFloat64(), longBody.ToFloat64()))
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
	i.Response(w, i.svc.AddInvProduct(bizID, invID, i.ProductIDField(r), i.Count(r), i.Available(r)))
}

func (i *InventoryHandler) InventoryUpdateProduct(w http.ResponseWriter, r *http.Request) {
	bizID := i.ClaimGetBusinessID(r)
	invProdID, err := i.GetInvProductID(r)
	if err != nil {
		i.ResponseBadRequest(w)
		return
	}
	i.Response(w, i.svc.UpdateInvProduct(bizID, invProdID, i.Count(r), i.Available(r)))
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
