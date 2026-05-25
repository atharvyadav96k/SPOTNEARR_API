package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
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

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	user, err := ParseBody[models.User](r)
	if err != nil {
		u.ResponseBadRequest(w)
	}
	res := u.GetUserService().RegisterUser(user)
	u.Response(w, res)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := ParseBody[models.User](r)
	if err != nil {
		u.ResponseBadRequest(w)
	}
	res := u.GetUserService().Login(*user.Email, *user.PasswordHash)
	u.Response(w, res)
}

func (u *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	u.ResponseOK(w)
}

func (u *UserHandler) BanUser(w http.ResponseWriter, r *http.Request) {
	u.ResponseOK(w)
}

func (u *UserHandler) UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	u.ResponseOK(w)
}
