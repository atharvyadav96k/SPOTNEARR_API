package handlers

import "net/http"

type ReviewHandler struct {
	BaseHandler
}

func NewReviewHandler() *ReviewHandler {
	return &ReviewHandler{}
}

func (rev *ReviewHandler) ReviewSpotlight(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewOffer(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewBusiness(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewProduct(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewSpotlightDelete(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewOfferDelete(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewBusinessDelete(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewProductDelete(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewSpotlightUpdate(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewOfferUpdate(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewBusinessUpdate(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewProductUpdate(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewGetBySpotlight(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewGetByOffer(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewGetByBusiness(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}

func (rev *ReviewHandler) ReviewGetByProduct(w http.ResponseWriter, r *http.Request) {
	rev.ResponseOK(w)
}
