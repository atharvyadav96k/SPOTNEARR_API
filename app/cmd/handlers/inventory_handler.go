package handlers

import "net/http"

type InventoryHandler struct {
	BaseHandler
}

func NewInventoryHandler() *InventoryHandler {
	return &InventoryHandler{}
}

func (i *InventoryHandler) InventoryCreate(w http.ResponseWriter, r *http.Request) {
	i.ResponseOK(w)
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
