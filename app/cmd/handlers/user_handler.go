package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type UserHandler struct {
	BaseHandler
}

func NewUserHandler(services *services.Services) *UserHandler {
	return &UserHandler{
		BaseHandler: *NewBaserHandler(services),
	}
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userId := u.ClaimGetUserId(r)
	if userId == 0 {
		u.ResponseBadRequest(w)
		return
	}
	res := u.GetUserService().GetUserProfile(userId)
	u.Response(w, res)
}

func (u *UserHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	u.ResponseOK(w)
}
