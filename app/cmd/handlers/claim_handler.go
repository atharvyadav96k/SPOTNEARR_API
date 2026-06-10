package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type ClaimHandler struct {
	BaseHandler
}

func NewClaimHandler(services *services.Services) *ClaimHandler {
	return &ClaimHandler{BaseHandler: *NewBaserHandler(services)}
}

func (c *ClaimHandler) GetUserClaims(w http.ResponseWriter, r *http.Request) {
	userID := c.ClaimGetUserId(r)
	if userID == 0 {
		c.ResponseBadRequest(w)
		return
	}
	c.Response(w, c.GetClaimService().GetUserClaims(userID))
}

func (c *ClaimHandler) ClaimProduct(w http.ResponseWriter, r *http.Request) {
	userID := c.ClaimGetUserId(r)
	if userID == 0 {
		c.ResponseBadRequest(w)
		return
	}
	invProductID, err := c.GetInvProductId(r)
	if err != nil {
		c.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	c.Response(w, c.GetClaimService().ClaimProduct(userID, invProductID))
}

func (c *ClaimHandler) ClaimRemove(w http.ResponseWriter, r *http.Request) {
	userID := c.ClaimGetUserId(r)
	if userID == 0 {
		c.ResponseBadRequest(w)
		return
	}
	claimID, err := c.GetClaimId(r)
	if err != nil {
		c.ResponseBadRequestWithMessage(w, err.Error())
		return
	}
	c.Response(w, c.GetClaimService().RemoveClaim(userID, claimID))
}
