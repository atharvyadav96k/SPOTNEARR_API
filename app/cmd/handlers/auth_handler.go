package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/atharvyadav96k/SPOTNEARR_API/auth"
	"github.com/atharvyadav96k/SPOTNEARR_API/config"
	"github.com/atharvyadav96k/SPOTNEARR_API/dtos"
	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type AuthHandler struct {
	BaseHandler
}

func NewAuthHandler(services *services.Services) *AuthHandler {
	return &AuthHandler{BaseHandler: *NewBaserHandler(services)}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var dto dtos.RegisterRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		a.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	a.Response(w, a.GetUserService().RegisterUser(dto.ToModel()))
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dto dtos.LoginRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		a.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	a.Response(w, a.GetUserService().Login(dto.Email, dto.Password))
}

func (a *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	var dto dtos.SessionRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		a.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	a.Response(w, a.GetUserService().SessionNotification(dto.Email))
}

func (a *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var dto dtos.ResetPasswordRequest
	if err := parseAndValidateBody(r, &dto); err != nil {
		a.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	a.Response(w, a.GetUserService().UpdatePassword(dto.Password, a.QuerySession(r)))
}

func (a *AuthHandler) Auth(w http.ResponseWriter, r *http.Request) {
	log.Default().Println("User Id: ", a.ClaimGetUserId(r))
	log.Default().Println("Business Id: ", a.ClaimGetBusinessId(r))
	a.ResponseOK(w)
}

func (a *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	userTokens, err := ParseBody[auth.TokenResponse](r)
	if err != nil || strings.TrimSpace(userTokens.RefreshToken) == "" {
		a.ResponseBadRequest(w)
		return
	}
	claims, err := auth.ValidateToken(userTokens.RefreshToken, config.C.JWTSecret, auth.TypeRefreshToken)
	if err != nil {
		a.ResponseBadRequest(w)
		return
	}
	a.Response(w, a.GetUserService().Refresh(*claims, userTokens.RefreshToken))
}

func (a *AuthHandler) LogoutFromAllDevices(w http.ResponseWriter, r *http.Request) {
	userId := a.ClaimGetUserId(r)
	if userId == 0 {
		a.ResponseBadRequest(w)
		return
	}
	a.Response(w, a.GetUserService().DismissRefreshToken(userId))
}
