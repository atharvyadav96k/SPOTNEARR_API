package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/spotnearr/vendor-svc/services"
)

type AuthHandler struct {
	BaseHandler
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto services.RegisterBusinessDTO
	if err := parseAndValidateBody(r, &dto); err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.svc.Register(&dto))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto services.LoginBusinessDTO
	if err := parseAndValidateBody(r, &dto); err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.svc.Login(&dto))
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var dto services.RefreshDTO
	if err := parseAndValidateBody(r, &dto); err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.svc.Refresh(&dto))
}
