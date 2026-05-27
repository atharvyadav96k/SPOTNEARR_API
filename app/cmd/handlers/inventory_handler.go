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
	i.ResponseOK(w)
}

func (i *InventoryHandler) InventoryDelete(w http.ResponseWriter, r *http.Request) {
	i.ResponseOK(w)
}

func (i *InventoryHandler) InventoryGetProducts(w http.ResponseWriter, r *http.Request) {
	i.ResponseOK(w)
}

func (i *InventoryHandler) InventoryAddProduct(w http.ResponseWriter, r *http.Request) {
	i.ResponseOK(w)
}

func (i *InventoryHandler) InventoryRemoveProduct(w http.ResponseWriter, r *http.Request) {
	i.ResponseOK(w)
}
