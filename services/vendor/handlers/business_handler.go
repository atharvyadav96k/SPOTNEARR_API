package handlers

import (
	"net/http"

	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
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
	var dto pkgdtos.RegisterBusiness
	if err := parseAndValidateBody(r, &dto); err != nil {
		b.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	b.Response(w, b.svc.RegisterBusiness(&dto, userID))
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
	var dto pkgdtos.UpdateBusiness
	if err := parseAndValidateBody(r, &dto); err != nil {
		b.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	b.Response(w, b.svc.UpdateBusiness(bizID, &dto))
}
