package handlers

import "net/http"

type OfferHandler struct {
	BaseHandler
}

func NewOfferHandler() *OfferHandler {
	return &OfferHandler{}
}

func (o *OfferHandler) OfferGet(w http.ResponseWriter, r *http.Request) {
	o.ResponseOK(w)
}

func (o *OfferHandler) NearByOffers(w http.ResponseWriter, r *http.Request) {
	o.ResponseOK(w)
}

func (o *OfferHandler) OfferAdd(w http.ResponseWriter, r *http.Request) {
	o.ResponseOK(w)
}

func (o *OfferHandler) OfferUpdate(w http.ResponseWriter, r *http.Request) {
	o.ResponseOK(w)
}

func (o *OfferHandler) OfferDelete(w http.ResponseWriter, r *http.Request) {
	o.ResponseOK(w)
}
