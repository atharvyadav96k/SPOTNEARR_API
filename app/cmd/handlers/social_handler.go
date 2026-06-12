package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type SocialHandler struct {
	BaseHandler
}

func NewSocialHandler(svcs *services.Services) *SocialHandler {
	return &SocialHandler{BaseHandler: *NewBaserHandler(svcs)}
}

func (h *SocialHandler) getSocialService() *services.SocialService {
	return h.services.SocialService
}

func (h *SocialHandler) FollowBusiness(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	bizID, err := h.GetBizId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().FollowBusiness(userID, bizID))
}

func (h *SocialHandler) UnfollowBusiness(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	bizID, err := h.GetBizId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().UnfollowBusiness(userID, bizID))
}

func (h *SocialHandler) LikeProduct(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	invProductID, err := h.GetInvProductId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().LikeProduct(userID, invProductID))
}

func (h *SocialHandler) UnlikeProduct(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	invProductID, err := h.GetInvProductId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().UnlikeProduct(userID, invProductID))
}

func (h *SocialHandler) SaveProduct(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	invProductID, err := h.GetInvProductId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().SaveProduct(userID, invProductID))
}

func (h *SocialHandler) UnsaveProduct(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	invProductID, err := h.GetInvProductId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().UnsaveProduct(userID, invProductID))
}

func (h *SocialHandler) LikeSpotlight(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	spotlightID, err := h.GetSpotlightId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().LikeSpotlight(userID, spotlightID))
}

func (h *SocialHandler) UnlikeSpotlight(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	spotlightID, err := h.GetSpotlightId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().UnlikeSpotlight(userID, spotlightID))
}

func (h *SocialHandler) SaveSpotlight(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	spotlightID, err := h.GetSpotlightId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().SaveSpotlight(userID, spotlightID))
}

func (h *SocialHandler) UnsaveSpotlight(w http.ResponseWriter, r *http.Request) {
	userID := h.ClaimGetUserId(r)
	if userID == 0 {
		h.ResponseBadRequest(w)
		return
	}
	spotlightID, err := h.GetSpotlightId(r)
	if err != nil {
		h.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	h.Response(w, h.getSocialService().UnsaveSpotlight(userID, spotlightID))
}
