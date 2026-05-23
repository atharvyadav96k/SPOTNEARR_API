package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type UserHandler struct {
	BaseHandler
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (u *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	user, err := ParseBody[models.User](r)
	if err != nil {
		u.ResponseBadRequest(w)
	}
	res := u.userService.RegisterUser(user)
	u.Response(w, res)
}

func (u *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	user, err := ParseBody[models.User](r)
	if err != nil {
		u.ResponseBadRequest(w)
	}
	res := u.userService.Login(*user.Email, *user.PasswordHash)
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
