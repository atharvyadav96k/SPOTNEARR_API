package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/dtos"
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
	userId := b.ClaimGetUserId(r)
	if userId == 0 {
		b.ResponseBadRequest(w)
		return
	}
	business, err := ParseBody[dtos.Business](r)
	if err != nil && business == nil {
		b.ResponseBadRequest(w)
		return
	}

	res := b.GetBizService().RegisterBusiness(business, userId)
	b.Response(w, res)
}

func (b *BusinessHandler) AddUserToBusiness(w http.ResponseWriter, r *http.Request) {
	// user, err := ParseBody[dtos.User](r)
	// if err != nil {
	// 	b.ResponseBadRequest(w)
	// 	return
	// }
	// res := b.GetUserService().AddUser()
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
	bizId := b.ClaimGetBusinessId(r)
	if bizId == 0 {
		b.ResponseBadRequest(w)
		return
	}
	name := b.Name(r)
	desc := b.Desc(r)
	res := b.GetBizService().UpdateBusinessById(bizId, bizId, name, desc)
	b.Response(w, res)
}

func (b *BusinessHandler) BusinessDelete(w http.ResponseWriter, r *http.Request) {
	b.ResponseOK(w)
}

func (b *BusinessHandler) BusinessInventories(w http.ResponseWriter, r *http.Request) {
	id := b.ClaimGetBusinessId(r)
	if id == 0 {
		b.ResponseBadRequest(w)
		return
	}
	res := b.GetInvService().GetBusinessInventory(id)
	b.Response(w, res)
}
