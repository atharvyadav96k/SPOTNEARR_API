package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type InventoryHandler struct {
	BaseHandler
}

func NewInventoryHandler(services *services.Services) *InventoryHandler {
	return &InventoryHandler{
		BaseHandler: *NewBaserHandler(services),
	}
}

func (i *InventoryHandler) InventoryCreate(w http.ResponseWriter, r *http.Request) {
	id := i.ClaimGetBusinessId(r)
	if id == 0 {
		i.ResponseBadRequest(w)
		return
	}
	store, err := ParseBody[models.Store](r)
	if err != nil || store == nil {
		i.ResponseBadRequest(w)
		return
	}
	res := i.GetInvService().CreateInventory(id, store)
	i.Response(w, res)
}

func (i *InventoryHandler) InventoryUpdate(w http.ResponseWriter, r *http.Request) {
	businessId := i.ClaimGetBusinessId(r)
	if businessId == 0 {
		i.ResponseBadRequest(w)
		return
	}
	storeId, err := i.GetInventoryId(r)
	if err != nil || storeId != 0 {
		i.ResponseBadRequest(w)
		return
	}
	name := i.Name(r)
	address := i.Address(r)
	lat := i.Lat(r).ToFloat64()
	long := i.Long(r).ToFloat64()

	res := i.GetInvService().UpdateInventory(businessId, storeId, name, address, lat, long)
	i.Response(w, res)
}

func (i *InventoryHandler) InventoryDelete(w http.ResponseWriter, r *http.Request) {
	i.ResponseOK(w)
}

func (i *InventoryHandler) InventoryGetProducts(w http.ResponseWriter, r *http.Request) {
	businessID := i.ClaimGetBusinessId(r)
	invId, err := i.GetInventoryId(r)
	if err != nil {
		i.ResponseBadRequest(w)
		return
	}
	res := i.GetInvService().GetInvProducts(businessID, invId)
	i.Response(w, res)
}

func (i *InventoryHandler) InventoryAddProduct(w http.ResponseWriter, r *http.Request) {
	businessID := i.ClaimGetBusinessId(r)
	invID, err := i.GetInventoryId(r)
	if err != nil {
		i.ResponseBadRequest(w)
		return
	}
	res := i.GetInvService().AddInvProduct(businessID, invID, i.ProductID(r), i.Count(r), i.Available(r))
	i.Response(w, res)
}

func (i *InventoryHandler) InventoryRemoveProduct(w http.ResponseWriter, r *http.Request) {
	i.ResponseOK(w)
}
