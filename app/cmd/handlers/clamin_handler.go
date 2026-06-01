package handlers

import (
	"net/http"

	"github.com/atharvyadav96k/SPOTNEARR_API/services"
)

type ClaimHandler struct {
	BaseHandler
}

func NewClaimHandler(services *services.Services) *ClaimHandler {
	return &ClaimHandler{
		BaseHandler: *NewBaserHandler(services),
	}
}

func (c *ClaimHandler) ClaimProduct(w http.ResponseWriter, r *http.Request) {
	userID := c.ClaimGetUserId(r)
	if userID == 0 {
		c.ResponseBadRequest(w)
		return
	}
	invProductID, err := c.GetInvProductId(r)
	if err != nil || invProductID == 0 {
		c.ResponseBadRequestWithMessage(w, "Invalid inventory product ID")
		return
	}
	res := c.GetClaimService().ClaimProduct(userID, invProductID)
	c.Response(w, res)
}

func (c *ClaimHandler) GetUserClaims(w http.ResponseWriter, r *http.Request) {
	userID := c.ClaimGetUserId(r)
	if userID == 0 {
		c.ResponseBadRequest(w)
		return
	}
	res := c.GetClaimService().GetUserClaims(userID)
	c.Response(w, res)
}

func (c *ClaimHandler) ClaimRemove(w http.ResponseWriter, r *http.Request) {
	userID := c.ClaimGetUserId(r)
	if userID == 0 {
		c.ResponseBadRequest(w)
		return
	}
	claimID, err := c.GetClaimId(r)
	if err != nil || claimID == 0 {
		c.ResponseBadRequestWithMessage(w, "Invalid claim ID")
		return
	}
	res := c.GetClaimService().RemoveClaim(userID, claimID)
	c.Response(w, res)
}
