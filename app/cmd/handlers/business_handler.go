package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type BusinessHandler struct {
	BaseHandler
}

func NewBusinessHandler(services *services.Services) *BusinessHandler {
	return &BusinessHandler{
		BaseHandler: *NewBaserHandler(services),
	}
}

func (b *BusinessHandler) BusinessRegister(w http.ResponseWriter, r *http.Request) {
	business, err := ParseBody[models.Business](r)
	if err != nil {
		b.ResponseBadRequest(w)
		return
	}
	res := b.GetBizService().RegisterBusiness(business)
	b.Response(w, res)
}

func (b *BusinessHandler) BusinessProfile(w http.ResponseWriter, r *http.Request) {
	id, err := b.GetBusinessId(r)
	if err != nil {
		b.ResponseBadRequest(w)
		return
	}
	res := b.GetBizService().GetBusinessById(id)
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
	res := b.GetBizService().UpdateBusinessById(id, business)
	b.Response(w, res)
}

func (b *BusinessHandler) BusinessDelete(w http.ResponseWriter, r *http.Request) {
	b.ResponseOK(w)
}

func (b *BusinessHandler) BusinessInventories(w http.ResponseWriter, r *http.Request) {
	id, err := b.GetBusinessId(r)
	if err != nil {
		b.ResponseBadRequest(w)
	}
	res := b.GetInvService().GetBusinessInventory(id)
	b.Response(w, res)
}
