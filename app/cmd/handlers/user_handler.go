package handlers

import (
	"net/http"

	pkgdtos "github.com/atharvyadav96k/spotnearr/pkg/dtos"
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
	u.Response(w, u.GetUserService().GetUserProfile(userId))
}

func (u *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := u.ClaimGetUserId(r)
	if userID == 0 {
		u.ResponseBadRequest(w)
		return
	}
	var dto pkgdtos.UpdateProfileRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		u.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	u.Response(w, u.GetUserService().UpdateProfile(userID, dto))
}

func (u *UserHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	targetID, err := u.GetUserId(r)
	if err != nil {
		u.ResponseBadRequestWithMessage(w, "invalid user ID")
		return
	}
	u.Response(w, u.GetUserService().BanUser(targetID))
}
