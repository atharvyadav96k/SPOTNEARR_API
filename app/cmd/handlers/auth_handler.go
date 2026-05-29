package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/auth"
	"github.com/atharvyadav96k/SPOTNEARR_API/models"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
	"github.com/atharvyadav96k/SPOTNEARR_API/utils/request"
)

type AuthHandler struct {
	BaseHandler
}

func NewAuthHandler(services *services.Services) *AuthHandler {
	return &AuthHandler{
		BaseHandler: *NewBaserHandler(services),
	}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	email, err := a.Email(r)
	if err != nil {
		a.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	password, err := a.Password(r)
	if err != nil {
		a.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	phone, err := a.Phone(r)
	if err != nil {
		a.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	user := models.NewUser(a.Name(r), email, phone, password)
	res := a.GetUserService().RegisterUser(user)
	a.Response(w, res)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email, err := a.Email(r)
	if err != nil {
		a.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	password := request.GetVal(r, "password").ToString()
	if strings.TrimSpace(password) == "" {
		a.ResponseBadRequestWithMessage(w, "missing password field")
		return
	}
	res := a.GetUserService().Login(email, password)
	a.Response(w, res)
}

func (a *AuthHandler) ResetSession(w http.ResponseWriter, r *http.Request) {
	// email, err := a.Email(r)
	// if err != nil {
	// 	a.ResponseBadRequestWithMessage(w, err.Error())
	// 	return
	// }
	// a.GetUserService().
}

func (a *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {

}

func (a *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("User Id: ", a.ClaimGetUserId(r))
	log.Default().Println("Business Id: ", a.ClaimGetBusinessId(r))
	a.ResponseOK(w)
}

func (a *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	userTokens, err := ParseBody[auth.TokenResponse](r)
	if err != nil {
		a.ResponseBadRequest(w)
		return
	}
	if strings.TrimSpace(userTokens.RefreshToken) == "" {
		a.ResponseBadRequest(w)
		return
	}
	claims, err := auth.ValidateToken(userTokens.RefreshToken, "dummy", auth.TypeRefreshToken)
	if err != nil {
		a.ResponseBadRequest(w)
		return
	}
	res := a.GetUserService().Refresh(*claims, userTokens.RefreshToken)
	a.Response(w, res)
}

func (a *AuthHandler) LogoutFromAllDevices(w http.ResponseWriter, r *http.Request) {
	userId := a.ClaimGetUserId(r)
	if userId == 0 {
		log.Default().Println("User Id: ", userId)
		a.ResponseBadRequest(w)
		return
	}
	res := a.GetUserService().DismissRefreshToken(userId)
	a.Response(w, res)
}
