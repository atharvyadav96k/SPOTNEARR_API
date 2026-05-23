package handlers

import "net/http"

type ProductHandler struct {
	BaseHandler
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

func (p *ProductHandler) ProductAdd(w http.ResponseWriter, r *http.Request) {
	p.ResponseOK(w)
}

func (p *ProductHandler) ProductUpdate(w http.ResponseWriter, r *http.Request) {
	p.ResponseOK(w)
}

func (p *ProductHandler) ProductDelete(w http.ResponseWriter, r *http.Request) {
	p.ResponseOK(w)
}

func (p *ProductHandler) ProductGet(w http.ResponseWriter, r *http.Request) {
	p.ResponseOK(w)
}

func (p *ProductHandler) ProductOfBusiness(w http.ResponseWriter, r *http.Request) {
	p.ResponseOK(w)
}

func (p *ProductHandler) ProductNearBy(w http.ResponseWriter, r *http.Request) {
	p.ResponseOK(w)
}
