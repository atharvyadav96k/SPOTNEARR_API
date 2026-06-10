package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/dtos"
	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
)

type BusinessHandler struct {
	BaseHandler
	svc *services.BusinessService
}

func NewBusinessHandler(svc *services.BusinessService) *BusinessHandler {
	return &BusinessHandler{svc: svc}
}

func (b *BusinessHandler) BusinessRegister(w http.ResponseWriter, r *http.Request) {
	userID := b.ClaimGetUserID(r)
	if userID == 0 {
		b.ResponseBadRequest(w)
		return
	}
	dto, err := ParseBody[dtos.RegisterBusiness](r)
	if err != nil || dto == nil {
		b.ResponseBadRequest(w)
		return
	}
	b.Response(w, b.svc.RegisterBusiness(dto, userID))
}

func (b *BusinessHandler) BusinessProfile(w http.ResponseWriter, r *http.Request) {
	id, err := b.GetBizID(r)
	if err != nil {
		b.ResponseBadRequest(w)
		return
	}
	b.Response(w, b.svc.GetBusinessByID(id))
}

func (b *BusinessHandler) BusinessUpdate(w http.ResponseWriter, r *http.Request) {
	bizID := b.ClaimGetBusinessID(r)
	if bizID == 0 {
		b.ResponseBadRequest(w)
		return
	}
	dto, err := ParseBody[dtos.UpdateBusiness](r)
	if err != nil || dto == nil {
		b.ResponseBadRequest(w)
		return
	}
	b.Response(w, b.svc.UpdateBusiness(bizID, dto))
}
