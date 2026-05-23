package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type BusinessHandler struct {
	BaseHandler
	businessService *services.BusinessService
}

func NewBusinessHandler(businessService *services.BusinessService) *BusinessHandler {
	return &BusinessHandler{
		businessService: businessService,
	}
}

func (b *BusinessHandler) BusinessRegister(w http.ResponseWriter, r *http.Request) {
	business, err := ParseBody[models.Business](r)
	if err != nil {
		b.ResponseBadRequest(w)
		return
	}
	res := b.businessService.RegisterBusiness(business)
	b.Response(w, res)
}

func (b *BusinessHandler) BusinessLogin(w http.ResponseWriter, r *http.Request) {
	b.ResponseOK(w)
}

func (b *BusinessHandler) BusinessProfile(w http.ResponseWriter, r *http.Request) {
	id, err := b.GetBusinessId(r)
	if err != nil {
		b.ResponseBadRequest(w)
		return
	}
	res := b.businessService.GetBusinessById(id)
	b.Response(w, res)
}

func (b *BusinessHandler) BusinessUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := b.GetBusinessId(r)
	if err != nil {
		b.ResponseBadRequest(w)
		return
	}
	business, err := ParseBody[models.Business](r)
	if err != nil {
		b.ResponseBadRequest(w)
	}
	res := b.businessService.UpdateBusinessById(id, business)
	b.Response(w, res)
}

func (b *BusinessHandler) BusinessDelete(w http.ResponseWriter, r *http.Request) {
	b.ResponseOK(w)
}

func (b *BusinessHandler) BusinessBan(w http.ResponseWriter, r *http.Request) {
	b.ResponseOK(w)
}

func (b *BusinessHandler) BusinessInventory(w http.ResponseWriter, t *http.Request) {
	b.ResponseOK(w)
}
